package handlers

import (
	"errors"
	"net/http"

	"github.com/faissalmaulana/go-approve/cmd/dto"
	"github.com/faissalmaulana/go-approve/internal/service"
	"github.com/faissalmaulana/go-approve/internal/service/auth"
	"github.com/faissalmaulana/go-approve/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

// ====== REGISTER HANDLER ======
type RegisterHandler struct {
	auth            *auth.Auth
	validate        *validator.Validate
	sugaredErrorMsg *utils.SugaredErrorMessageValidator
}

func NewRegisterHandler(a *auth.Auth, val *validator.Validate, s *utils.SugaredErrorMessageValidator) *RegisterHandler {
	return &RegisterHandler{
		auth:            a,
		validate:        val,
		sugaredErrorMsg: s,
	}
}

func (r *RegisterHandler) HandleFunc(c *echo.Context) error {
	userPayload := new(dto.UserRegisterDTO)
	if err := c.Bind(userPayload); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
	}

	if err := r.validate.Struct(userPayload); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse(r.sugaredErrorMsg.TranslateValidationErrors(err)))
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
			return c.JSON(http.StatusConflict, utils.ErrorResponse(err.Error()))
		case errors.Is(err, service.ErrInvalidPayload):
			return c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		}
	}

	return c.JSON(http.StatusCreated, utils.SuccessResponse(map[string]string{"message": "Register Successfully", "access_token": accessToken}))
}

// ====== LOGIN HANDLER ======
type LoginHandler struct {
	auth            *auth.Auth
	validate        *validator.Validate
	sugaredErrorMsg *utils.SugaredErrorMessageValidator
}

func NewLoginHandler(a *auth.Auth, val *validator.Validate, s *utils.SugaredErrorMessageValidator) *LoginHandler {
	return &LoginHandler{
		auth:            a,
		validate:        val,
		sugaredErrorMsg: s,
	}
}

func (l *LoginHandler) HandleFunc(c *echo.Context) error {
	userPayload := new(dto.UserLoginDTO)
	if err := c.Bind(userPayload); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
	}

	if err := l.validate.Struct(userPayload); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse(l.sugaredErrorMsg.TranslateValidationErrors(err)))
	}

	accessToken, err := l.auth.Login(c.Request().Context(), userPayload.Email, userPayload.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidPayload):
			return c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		case errors.Is(err, service.ErrPasswordNotMatched), errors.Is(err, service.ErrUserNotFound):
			return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Password or Email is incorrect"))
		default:
			return c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		}
	}

	return c.JSON(http.StatusCreated, utils.SuccessResponse(map[string]string{"message": "Login Successfully", "access_token": accessToken}))
}
