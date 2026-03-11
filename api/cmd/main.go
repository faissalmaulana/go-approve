package main

import (
	"net/http"

	"github.com/faissalmaulana/go-approve/cmd/handlers"
	"github.com/faissalmaulana/go-approve/internal/db"
	"github.com/faissalmaulana/go-approve/internal/repository/user"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	fx.New(
		fx.Provide(
			NewEchoMux,
			NewHttpServer,
			db.New,
			fx.Annotate(user.New, fx.As(new(user.UserStorage))),
			handlers.NewHealthHandler,
			handlers.NewRegisterHandler,
			zap.NewProduction,
		),
		fx.Invoke(func(*gorm.DB) {}),
		fx.Invoke(func(*http.Server) {})).Run()
}
