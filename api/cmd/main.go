package main

import (
	"net/http"

	"github.com/faissalmaulana/go-approve/cmd/handlers"
	"github.com/faissalmaulana/go-approve/internal/db"
	"github.com/faissalmaulana/go-approve/internal/repository/user"
	"github.com/faissalmaulana/go-approve/internal/service/auth"
	"github.com/faissalmaulana/go-approve/internal/service/jwtfx"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	fx.New(
		jwtfx.Module,
		fx.Provide(
			NewEchoMux,
			NewHttpServer,
			auth.New,
			db.New,
			fx.Annotate(user.New, fx.As(new(user.UserStorage))),
			handlers.NewHealthHandler,
			handlers.NewRegisterHandler,
			zap.NewProduction,
		),
		fx.Invoke(func(*gorm.DB) {}),
		fx.Invoke(func(*http.Server) {})).Run()
}
