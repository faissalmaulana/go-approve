package main

import (
	"net/http"

	"github.com/faissalmaulana/go-approve/cmd/routes"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"go.uber.org/fx"
)

type EchoMuxParams struct {
	fx.In

	Health *routes.HealthHandler
}

func NewEchoMux(p EchoMuxParams) http.Handler {
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/health", p.Health.HandleFunc)

	return e
}
