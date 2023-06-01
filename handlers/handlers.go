package handlers

import (
	"html/template"

	"github.com/charmbracelet/log"

	"github.com/labstack/echo/v4"
)

func IndexHandler(c echo.Context) error {
	t, err := template.ParseFiles("views/base.html", "views/public/index.html")
	if err != nil {
		log.Error("failed to parse template", "err", err)
	}

	return t.Execute(c.Response().Writer, nil)
}
