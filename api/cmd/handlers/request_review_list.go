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

type GetRequestReviewHandler struct {
	validate             *validator.Validate
	sugaredErrorMsg      *utils.SugaredErrorMessageValidator
	requestReviewService *requestreview.RequestReviewService
	user                 *user.User
}

func NewGetRequestReviewHandler(
	v *validator.Validate,
	sgr *utils.SugaredErrorMessageValidator,
	rrs *requestreview.RequestReviewService,
	usr *user.User,
) *GetRequestReviewHandler {
	return &GetRequestReviewHandler{
		validate:             v,
		sugaredErrorMsg:      sgr,
		requestReviewService: rrs,
		user:                 usr,
	}
}

func (g *GetRequestReviewHandler) HandleFunc(c *echo.Context) error {
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

	var q = new(dto.GetRequestReviewDTO)
	if err := c.Bind(q); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
	}

	if err := g.validate.Struct(q); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse(g.sugaredErrorMsg.TranslateValidationErrors(err)))
	}

	as := q.As
	if as == "" {
		as = "received"
	}

	var status *utils.Status
	if q.Status != "" {
		var parsed utils.Status
		switch q.Status {
		case "pending":
			parsed = utils.Pending
		case "accepted":
			parsed = utils.Accepted
		case "rejected":
			parsed = utils.Rejected
		}
		status = &parsed
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

	type invitedUser struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Handler string `json:"handler"`
	}
	type responseRow struct {
		ID        string      `json:"id"`
		RoomID    string      `json:"room_id"`
		Status    string      `json:"status"`
		CreatedAt string      `json:"created_at"`
		User      invitedUser `json:"user"`
	}

	if as == "received" {
		items, err := g.requestReviewService.GetReceivedInvitations(
			c.Request().Context(),
			currentUser.ID,
			status,
			limit,
			offset,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		}

		resp := make([]responseRow, 0, len(items))
		for _, it := range items {
			resp = append(resp, responseRow{
				ID:        it.ID,
				RoomID:    it.ApprovalRoomId,
				Status:    it.Status,
				CreatedAt: it.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
				User: invitedUser{
					ID:      it.Requester.ID,
					Name:    it.Requester.Name,
					Handler: it.Requester.Handler,
				},
			})
		}

		return c.JSON(http.StatusOK, utils.SuccessResponse(resp))
	}

	items, err := g.requestReviewService.GetSentInvitations(
		c.Request().Context(),
		currentUser.ID,
		status,
		limit,
		offset,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
	}

	resp := make([]responseRow, 0, len(items))
	for _, it := range items {
		resp = append(resp, responseRow{
			ID:        it.ID,
			RoomID:    it.ApprovalRoomId,
			Status:    it.Status,
			CreatedAt: it.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			User: invitedUser{
				ID:      it.Invitee.ID,
				Name:    it.Invitee.Name,
				Handler: it.Invitee.Handler,
			},
		})
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse(resp))
}
