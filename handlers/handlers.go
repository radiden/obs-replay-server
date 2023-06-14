package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/andreykaipov/goobs/api/events"
	"github.com/charmbracelet/log"
	"github.com/goodsign/monday"
	"github.com/martinhoefling/goxkcdpwgen/xkcdpwgen"
	"github.com/radiden/obs-replay-server/models"
	"github.com/radiden/obs-replay-server/services"

	"github.com/labstack/echo/v4"
)

var sv *services.Services
var g = xkcdpwgen.NewGenerator()

func InitHandlers(svs *services.Services) {
	sv = svs
	g.SetNumWords(3)
	g.SetDelimiter("-")
	g.SetCapitalize(false)
}

func IndexHandler(c echo.Context) error {
	t, err := template.ParseFiles("views/base.html", "views/public/index.html")
	if err != nil {
		return err
	}

	return t.Execute(c.Response().Writer, nil)
}

func AuthHandler(c echo.Context) error {
	userName := c.FormValue("username")
	c.SetCookie(&http.Cookie{
		Name:  "sdvxreplay_session",
		Value: userName,
	})

	_, err := sv.DB.Queries.UserByName(c.Request().Context(), userName)
	if err == sql.ErrNoRows {
		_, err := sv.DB.Queries.CreateUser(c.Request().Context(), userName)
		if err != nil {
			log.Error("couldn't create new user", "err", err)
			return err
		}
	} else if err != nil {
		log.Error("couldn't query user", "err", err)
		return err
	}

	t, err := template.ParseFiles("views/base.html", "views/public/login_success.html")
	if err != nil {
		log.Error("couldn't parse template", "err", err)
		return err
	}

	return t.Execute(c.Response().Writer, nil)
}

func PanelHandler(c echo.Context) error {
	ctx := c.Request().Context()

	cookie, err := c.Cookie("sdvxreplay_session")
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	replays, err := sv.DB.Queries.ReplaysForUser(ctx, cookie.Value)
	if err != nil {
		log.Error("failed to get replays", "err", err)
	}

	t, err := template.ParseFiles("views/base.html", "views/public/panel.html")
	if err != nil {
		log.Error("failed to render panel template", "err", err)
		return err
	}

	return t.Execute(c.Response().Writer, map[string]any{
		"replays":  replays,
		"username": cookie.Value,
	})
}

func SaveHandler(c echo.Context) error {
	ctx := c.Request().Context()
	cookie, err := c.Cookie("sdvxreplay_session")
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	path, err := saveOBSReplay()
	if err != nil {
		return err
	}

	log.Info("saved file", "path", path)

	user, err := sv.DB.Queries.UserByName(ctx, cookie.Value)
	if err != nil {
		log.Error("couldn't query user in db", "user", cookie.Value, "err", err)
	}

	sv.DB.Queries.AddReplay(ctx, models.AddReplayParams{
		FilePath:     path,
		CreationTime: monday.Format(time.Now(), "Monday, 02 Jan 2006 15:04:05", monday.LocalePlPL),
		Owner:        user.ID,
	})

	return c.JSON(http.StatusOK, map[string]string{
		"path": path,
	})
}

func saveOBSReplay() (string, error) {
	_, err := sv.OBS.Outputs.SaveReplayBuffer()
	if err != nil {
		log.Error("failed saving replay buffer", "err", err)
		return "", err
	}

	// get saved replay path
	savedFilePath := getSavedPath()

	// open saved replay
	inputFile, err := os.Open(savedFilePath)
	if err != nil {
		log.Error("couldn't open source replay file", "err", err)
		return "", err
	}

	// get cwd
	currentDir, err := os.Getwd()
	if err != nil {
		log.Error("couldn't get cwd", "err", err)
		return "", err
	}

	// create output file
	fileName := g.GeneratePasswordString()
	destDir := filepath.Join(currentDir, "replays")
	destFileName := fmt.Sprintf("%s%s", fileName, filepath.Ext(savedFilePath))
	outputFile, err := os.Create(filepath.Join(destDir, destFileName))
	if err != nil {
		inputFile.Close()
		log.Error("couldn't create output replay file", "err", err)
		return "", err
	}
	// close output file after function runs
	defer outputFile.Close()

	// copy source to output
	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		log.Error("failed to copy replay file", "err", err)
		return "", err
	}

	// close input file
	inputFile.Close()
	// remove input file
	err = os.Remove(savedFilePath)
	if err != nil {
		log.Error("failed removing source replay file", "err", err)
	}

	return destFileName, nil
}

func getSavedPath() string {
	for {
		msg := <-sv.OBS.IncomingEvents
		switch m := msg.(type) {
		case *events.ReplayBufferSaved:
			return m.SavedReplayPath
		}
	}
}
