package main

import (
	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/radiden/obs-replay-server/handlers"
	"github.com/radiden/obs-replay-server/services"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed db/schema.sql
var initialSchema string

var sv *services.Services

func main() {
	sv = services.InitServices(initialSchema)
	handlers.InitHandlers(sv)

	isRecording, err := sv.OBS.Outputs.GetReplayBufferStatus()
	if err != nil {
		log.Fatal("couldn't get replay buffer status", "error", err)
	}
	if !isRecording.OutputActive {
		sv.OBS.Outputs.StartReplayBuffer()
	}

	// start echo
	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogError:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if err == nil {
				log.Info("[request]", "uri", v.URI, "status", v.Status)
			} else {
				log.Error("[request]", "uri", v.URI, "status", v.Status, "err", v.Error)
			}
			return nil
		},
	}))
	e.Debug = true
	e.HideBanner = true
	e.Static("/replays", "replays")
	e.GET("/", handlers.IndexHandler)
	e.POST("/auth", handlers.AuthHandler)
	e.GET("/panel", handlers.PanelHandler)
	e.GET("/save", handlers.SaveHandler)
	e.Logger.Fatal(e.Start(":80"))
}
