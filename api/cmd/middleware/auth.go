package middleware

import (
	"net/http"
	"strings"

	"github.com/faissalmaulana/go-approve/internal/repository/blocklisttoken"
	"github.com/faissalmaulana/go-approve/internal/utils"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type AuthMiddleware struct {
	blocklist blocklisttoken.BlocklistStorage
}

func NewAuthMiddleware(b blocklisttoken.BlocklistStorage) *AuthMiddleware {
	return &AuthMiddleware{
		blocklist: b,
	}
}

func (m *AuthMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("authorization header is required"))
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		_, err := m.blocklist.FindByToken(c.Request().Context(), tokenString)
		if err == nil {
			return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("token has been revoked"))
		}

		if err != gorm.ErrRecordNotFound {
			return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("failed to validate token"))
		}

		c.Set("token", tokenString)

		return next(c)
	}
}
