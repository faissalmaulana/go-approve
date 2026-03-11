package handlers

import (
	"net/http"

	"github.com/faissalmaulana/go-approve/internal/service"
	"github.com/faissalmaulana/go-approve/internal/service/user"
	"github.com/faissalmaulana/go-approve/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

type UserProfileHandler struct {
	user *user.User
}

func NewUserProfileHandler(u *user.User) *UserProfileHandler {
	return &UserProfileHandler{
		user: u,
	}
}

func (u *UserProfileHandler) HandleFunc(c *echo.Context) error {
	token, err := echo.ContextGet[*jwt.Token](c, "user")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse(err.Error()))
	}

	claims := token.Claims.(*jwt.RegisteredClaims)
	if claims.Subject == "" {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse(service.ErrSubIsEmpty.Error()))
	}

	currentUser, err := u.user.GetCurrentUser(c.Request().Context(), claims.Subject)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			return c.JSON(http.StatusNotFound, utils.ErrorResponse(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		}
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse(currentUser))
}

type LogoutHandler struct {
	user *user.User
}

func NewLogoutHandler(u *user.User) *LogoutHandler {
	return &LogoutHandler{
		user: u,
	}
}

func (l *LogoutHandler) HandleFunc(c *echo.Context) error {
	token, ok := c.Get("token").(string)
	if !ok || token == "" {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("token not found"))
	}

	err := l.user.Logout(c.Request().Context(), token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse(map[string]string{"message": "Logout successfully"}))
}
