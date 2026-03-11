package handlers

import (
	"net/http"

	"github.com/faissalmaulana/go-approve/cmd/dto"
	"github.com/faissalmaulana/go-approve/internal/repository/user"
	"github.com/labstack/echo/v5"
	"golang.org/x/crypto/bcrypt"
)

// ====== REGISTER HANDLER ======
type RegisterHandler struct {
	UserStore user.UserStorage
}

func NewRegisterHandler(u user.UserStorage) *RegisterHandler {
	return &RegisterHandler{
		UserStore: u,
	}
}

func (r *RegisterHandler) HandleFunc(c *echo.Context) error {
	userPayload := new(dto.UserDTO)
	if err := c.Bind(userPayload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPayload.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err := r.UserStore.Create(
		c.Request().Context(),
		userPayload.Email,
		userPayload.Name,
		string(hashedPassword),
		userPayload.Handler,
	); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Register Successfully"})
}

// ====== LOGIN HANDLER ======
type LoginHandler struct {
}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{}
}

func (*LoginHandler) HandleFunc(c *echo.Context) error {

	return c.String(http.StatusOK, "OK")
}
