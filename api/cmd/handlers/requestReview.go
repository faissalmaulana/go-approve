package handlers

import (
	"net/http"

	"github.com/faissalmaulana/go-approve/cmd/dto"
	"github.com/faissalmaulana/go-approve/internal/service"
	requestreview "github.com/faissalmaulana/go-approve/internal/service/requestReview"
	"github.com/faissalmaulana/go-approve/internal/service/user"
	"github.com/faissalmaulana/go-approve/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

type ConfirmRequestReviewHandler struct {
	validate             *validator.Validate
	sugaredErrorMsg      *utils.SugaredErrorMessageValidator
	requestReviewService *requestreview.RequestReviewService
	user                 *user.User
}

func NewConfirmRequestReviewHandler(
	v *validator.Validate,
	sgr *utils.SugaredErrorMessageValidator,
	rrs *requestreview.RequestReviewService,
	usr *user.User,
) *ConfirmRequestReviewHandler {
	return &ConfirmRequestReviewHandler{
		validate:             v,
		sugaredErrorMsg:      sgr,
		requestReviewService: rrs,
		user:                 usr,
	}
}

func (c *ConfirmRequestReviewHandler) HandleFunc(cc *echo.Context) error {
	token, err := echo.ContextGet[*jwt.Token](cc, "user")
	if err != nil {
		return cc.JSON(http.StatusUnauthorized, utils.ErrorResponse(err.Error()))
	}

	claims := token.Claims.(*jwt.RegisteredClaims)
	if claims.Subject == "" {
		return cc.JSON(http.StatusUnauthorized, utils.ErrorResponse(service.ErrSubIsEmpty.Error()))
	}

	requestReviewId := cc.Param("id")

	currentUser, err := c.user.GetCurrentUser(cc.Request().Context(), claims.Subject)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			return cc.JSON(http.StatusNotFound, utils.ErrorResponse(err.Error()))
		default:
			return cc.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		}
	}

	var confirmDTO = new(dto.ConfirmReviewRequestDTO)
	if err := cc.Bind(confirmDTO); err != nil {
		return cc.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
	}

	if err := c.validate.Struct(confirmDTO); err != nil {
		return cc.JSON(http.StatusBadRequest, utils.ErrorResponse(c.sugaredErrorMsg.TranslateValidationErrors(err)))
	}

	status := utils.Accepted
	if confirmDTO.Status == "rejected" {
		status = utils.Rejected
	}

	err = c.requestReviewService.ConfirmRequestRevieWithInsertApprover(
		cc.Request().Context(),
		requestReviewId,
		currentUser.ID,
		status,
	)
	if err != nil {
		return cc.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
	}

	return cc.JSON(http.StatusOK, utils.SuccessResponse(map[string]string{
		"message": "Success confirm review request",
	}))
}
