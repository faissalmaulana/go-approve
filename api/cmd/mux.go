package main

import (
	"net/http"
	"os"

	"github.com/faissalmaulana/go-approve/cmd/handlers"
	authMiddleware "github.com/faissalmaulana/go-approve/cmd/middleware"
	"github.com/faissalmaulana/go-approve/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	echomiddleware "github.com/labstack/echo/v5/middleware"
	"go.uber.org/fx"
)

type EchoMuxParams struct {
	fx.In

	Health   *handlers.HealthHandler
	Register *handlers.RegisterHandler
	Login    *handlers.LoginHandler
	Profile  *handlers.UserProfileHandler
	Logout   *handlers.LogoutHandler
	Auth     *authMiddleware.AuthMiddleware
}

func NewEchoMux(p EchoMuxParams) http.Handler {
	e := echo.New()
	e.Use(echomiddleware.RequestLogger())
	e.Use(echomiddleware.Recover())

	e.GET("/health", p.Health.HandleFunc)

	e.POST("/auth/sign-up", p.Register.HandleFunc)
	e.POST("/auth/sign-in", p.Login.HandleFunc)

	// protected routes
	r := e.Group("")

	config := echojwt.Config{
		NewClaimsFunc: func(c *echo.Context) jwt.Claims {
			return new(jwt.RegisteredClaims)
		},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
		ErrorHandler: func(c *echo.Context, err error) error {
			return c.JSON(http.StatusUnauthorized, utils.ErrorResponse(err.Error()))
		},
	}

	r.Use(echojwt.WithConfig(config))
	r.Use(p.Auth.Authenticate)

	r.GET("/profile", p.Profile.HandleFunc)
	r.POST("/logout", p.Logout.HandleFunc)

	return e
}
