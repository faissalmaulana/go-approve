package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type HealthHandler struct {
	// fields of dependencies
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (*HealthHandler) HandleFunc(c *echo.Context) error {
	return c.String(http.StatusOK, "ok")
}
