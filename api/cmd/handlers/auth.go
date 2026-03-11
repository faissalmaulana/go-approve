package handlers

import (
	"errors"
	"net/http"

	"github.com/faissalmaulana/go-approve/cmd/dto"
	"github.com/faissalmaulana/go-approve/internal/service"
	"github.com/faissalmaulana/go-approve/internal/service/auth"
	"github.com/labstack/echo/v5"
)

// ====== REGISTER HANDLER ======
type RegisterHandler struct {
	auth *auth.Auth
}

func NewRegisterHandler(a *auth.Auth) *RegisterHandler {
	return &RegisterHandler{
		auth: a,
	}
}

func (r *RegisterHandler) HandleFunc(c *echo.Context) error {
	userPayload := new(dto.UserDTO)
	if err := c.Bind(userPayload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	accessToken, err := r.auth.Register(
		c.Request().Context(),
		userPayload.Email,
		userPayload.Name,
		userPayload.Password,
		userPayload.Handler,
	)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrDuplicatedUser):
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		case errors.Is(err, service.ErrInvalidPayload):
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Register Successfully", "access_token": accessToken})
}

// ====== LOGIN HANDLER ======
type LoginHandler struct {
	auth *auth.Auth
}

func NewLoginHandler(a *auth.Auth) *LoginHandler {
	return &LoginHandler{
		auth: a,
	}
}

func (l *LoginHandler) HandleFunc(c *echo.Context) error {
	userPayload := new(dto.UserLoginDTO)
	if err := c.Bind(userPayload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	accessToken, err := l.auth.Login(c.Request().Context(), userPayload.Email, userPayload.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidPayload):
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrPasswordNotMatched), errors.Is(err, service.ErrUserNotFound):
			return echo.NewHTTPError(http.StatusBadRequest, "Password or Email is incorrect")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Login Successfully", "access_token": accessToken})
}
