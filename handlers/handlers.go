package handlers

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/radiden/obs-replay-server/services"

	"github.com/labstack/echo/v4"
)

var sv *services.Services

func InitHandlers(svs *services.Services) {
	sv = svs
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
			return err
		}
	} else if err != nil {
		return err
	}

	t, err := template.ParseFiles("views/base.html", "views/public/login_success.html")
	if err != nil {
		return err
	}

	return t.Execute(c.Response().Writer, nil)
}

func PanelHandler(c echo.Context) error {
	_, err := c.Cookie("sdvxreplay_session")
	if err != nil {
		return c.Redirect(http.StatusBadRequest, "/")
	}
	t, err := template.ParseFiles("views/base.html", "views/public/panel.html")
	if err != nil {
		return err
	}

	return t.Execute(c.Response().Writer, nil)
}
