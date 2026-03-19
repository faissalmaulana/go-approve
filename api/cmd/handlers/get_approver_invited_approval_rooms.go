package handlers

import (
	"net/http"

	"github.com/faissalmaulana/go-approve/cmd/dto"
	"github.com/faissalmaulana/go-approve/internal/service"
	approvalroom "github.com/faissalmaulana/go-approve/internal/service/approvalRoom"
	"github.com/faissalmaulana/go-approve/internal/service/user"
	"github.com/faissalmaulana/go-approve/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

type GetApproverInvitedApprovalRoomsHandler struct {
	validate        *validator.Validate
	sugaredErrorMsg *utils.SugaredErrorMessageValidator

	approvalRoomService *approvalroom.ApprovalRoomService
	user                *user.User
}

func NewGetApproverInvitedApprovalRoomsHandler(
	v *validator.Validate,
	sgr *utils.SugaredErrorMessageValidator,
	ars *approvalroom.ApprovalRoomService,
	usr *user.User,
) *GetApproverInvitedApprovalRoomsHandler {
	return &GetApproverInvitedApprovalRoomsHandler{
		validate:             v,
		sugaredErrorMsg:      sgr,
		approvalRoomService: ars,
		user:                 usr,
	}
}

func (g *GetApproverInvitedApprovalRoomsHandler) HandleFunc(c *echo.Context) error {
	token, err := echo.ContextGet[*jwt.Token](c, "user")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse(err.Error()))
	}

	claims := token.Claims.(*jwt.RegisteredClaims)
	if claims.Subject == "" {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse(service.ErrSubIsEmpty.Error()))
	}

	currentUser, err := g.user.GetCurrentUser(c.Request().Context(), claims.Subject)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			return c.JSON(http.StatusNotFound, utils.ErrorResponse(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		}
	}

	var q = new(dto.GetApproverInvitedApprovalRoomsDTO)
	if err := c.Bind(q); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
	}

	if err := g.validate.Struct(q); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse(g.sugaredErrorMsg.TranslateValidationErrors(err)))
	}

	sortField := q.Sort
	if sortField == "" {
		sortField = "created_at"
	}

	orderDir := q.Order
	if orderDir == "" {
		orderDir = "desc"
	}

	limit := q.Limit
	if limit == 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := q.Offset
	if offset < 0 {
		offset = 0
	}

	rooms, err := g.approvalRoomService.GetApproverInvitedApprovalRooms(
		c.Request().Context(),
		currentUser.ID,
		sortField,
		orderDir,
		limit,
		offset,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse(rooms))
}

