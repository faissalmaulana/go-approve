package main

import (
	"net/http"

	"github.com/faissalmaulana/go-approve/cmd/handlers"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.Provide(
			NewEchoMux,
			NewHttpServer,
			handlers.NewHealthHandler,
			zap.NewProduction,
		),
		fx.Invoke(func(*http.Server) {})).Run()
}
