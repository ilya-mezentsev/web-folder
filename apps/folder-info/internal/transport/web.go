package transport

import (
	"errors"
	"github.com/ilya-mezentsev/folder-info/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogEcho "github.com/samber/slog-echo"
	"log/slog"
	"net/http"
)

type Error struct {
	Val string `json:"error"`
}

func Run(
	addr string,
	core services.Core,
	logger *slog.Logger,
) {
	e := echo.New()

	e.Use(slogEcho.New(logger))
	e.Use(middleware.Recover())

	e.GET("/roots", func(c echo.Context) error {
		return c.JSON(http.StatusOK, core.Roots())
	})
	e.GET("/folder", infoFetcher(core.Folder))
	e.GET("/file", infoFetcher(core.File))

	e.Logger.Fatal(e.Start(addr))
}

func infoFetcher[T any](fetch func(string) (T, error)) func(ctx echo.Context) error {
	return func(c echo.Context) error {
		info, err := fetch(c.QueryParam("path"))
		if err != nil {
			if errors.Is(err, services.ErrPathIsNotAllowed) {
				return c.NoContent(http.StatusBadRequest)
			}

			return c.JSON(http.StatusInternalServerError, Error{Val: err.Error()})
		}

		return c.JSON(http.StatusOK, info)
	}
}
