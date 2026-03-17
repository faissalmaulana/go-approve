package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/faissalmaulana/go-approve/cmd/dto"
	"github.com/faissalmaulana/go-approve/internal/model"
	"github.com/faissalmaulana/go-approve/internal/service"
	approvalroom "github.com/faissalmaulana/go-approve/internal/service/approvalRoom"
	contractapprovalroom "github.com/faissalmaulana/go-approve/internal/service/approvalRoom/contract"
	"github.com/faissalmaulana/go-approve/internal/service/user"
	"github.com/faissalmaulana/go-approve/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

type CreateApprovalRoomHandler struct {
	validate            *validator.Validate
	sugaredErrorMsg     *utils.SugaredErrorMessageValidator
	approvalRoomService *approvalroom.ApprovalRoomService
	user                *user.User
}

func NewCreateApprovalRoomHandler(
	v *validator.Validate,
	sgr *utils.SugaredErrorMessageValidator,
	ars *approvalroom.ApprovalRoomService,
	usr *user.User,
) *CreateApprovalRoomHandler {
	return &CreateApprovalRoomHandler{
		validate:            v,
		sugaredErrorMsg:     sgr,
		approvalRoomService: ars,
		user:                usr,
	}
}

func (a *CreateApprovalRoomHandler) HandleFunc(c *echo.Context) error {
	token, err := echo.ContextGet[*jwt.Token](c, "user")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse(err.Error()))
	}

	claims := token.Claims.(*jwt.RegisteredClaims)
	if claims.Subject == "" {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse(service.ErrSubIsEmpty.Error()))
	}

	var jsonCreateApprovalRoom = new(dto.CreateApprovalRoomDTO)

	createApprovalRoomPayload := c.FormValue("json_data")

	if err := json.Unmarshal([]byte(createApprovalRoomPayload), jsonCreateApprovalRoom); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
	}

	if err := a.validate.Struct(jsonCreateApprovalRoom); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse(a.sugaredErrorMsg.TranslateValidationErrors(err)))
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
	}

	files := form.File["documents"]
	fileMetadatas := make(map[string]model.FileMetadata, len(files))
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		}

		defer src.Close()

		generatedfilepath := filepath.Join("storage", "private", utils.GenerateRandomFilename(file.Filename))
		dst, err := os.Create(generatedfilepath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		}

		if _, err := io.Copy(dst, src); err != nil {
			return c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		}

		fileMetadatas[file.Filename] = model.FileMetadata{
			Link:     generatedfilepath,
			Filename: file.Filename,
			Size:     int(file.Size),
		}
	}

	currentUser, err := a.user.GetCurrentUser(c.Request().Context(), claims.Subject)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			return c.JSON(http.StatusNotFound, utils.ErrorResponse(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		}
	}

	// GET USERS APPROVALS
	approverIds, err := a.user.GetUserIdsOnly(c.Request().Context(), jsonCreateApprovalRoom.Approvers)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
	}

	newApprovalRoom := &contractapprovalroom.CreateApprovalRoom{
		Title:         jsonCreateApprovalRoom.Title,
		DueAt:         jsonCreateApprovalRoom.DueAt,
		SubmitterId:   currentUser.ID,
		InviteeIds:    approverIds,
		FileMetadatas: fileMetadatas,
	}

	// CREATE ROOM
	newApprovalRoomId, err := a.approvalRoomService.Create(c.Request().Context(), newApprovalRoom)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
	}

	return c.JSON(http.StatusCreated, utils.SuccessResponse(map[string]string{
		"message":              "Success Create Room",
		"new_approval_room_id": newApprovalRoomId,
	}))
}
