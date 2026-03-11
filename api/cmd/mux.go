package main

import (
	"net/http"

	"github.com/faissalmaulana/go-approve/cmd/handlers"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"go.uber.org/fx"
)

type EchoMuxParams struct {
	fx.In

	Health   *handlers.HealthHandler
	Register *handlers.RegisterHandler
	Login    *handlers.LoginHandler
}

func NewEchoMux(p EchoMuxParams) http.Handler {
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/health", p.Health.HandleFunc)
	e.POST("/signup", p.Register.HandleFunc)
	e.POST("/sign-in", p.Login.HandleFunc)

	return e
}
