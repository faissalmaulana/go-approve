package main

import (
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.Provide(
			NewEchoMux,
			NewHttpServer,
			zap.NewProduction,
		),
		fx.Invoke(func(*http.Server) {})).Run()
}
