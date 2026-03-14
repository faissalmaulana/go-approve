package approvalroom

import (
	"context"
	"strings"
	"time"

	"github.com/faissalmaulana/go-approve/internal/constant"
	"github.com/faissalmaulana/go-approve/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApprovalRoomStorage interface {
	Create(ctx context.Context,
		title string,
		filepaths []string,
		dueAt time.Time,
		submitterId string,
	) error
	CreateWithTx(
		title string,
		filepaths []string,
		dueAt time.Time,
		submitterId string,
		approvalRoomId *string,
	) func(ctx context.Context, tx *gorm.DB) error
}

type ApprovalRoomRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) ApprovalRoomStorage {
	return &ApprovalRoomRepository{
		db: db,
	}
}

func (ar *ApprovalRoomRepository) Create(
	ctx context.Context,
	title string,
	filepaths []string,
	dueAt time.Time,
	submitterId string,
) error {

	ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
	defer cancel()

	newRoom := &model.ApprovalRoom{
		Title:       title,
		Filepaths:   strings.Join(filepaths, ";"),
		DueAt:       dueAt,
		SubmitterId: submitterId,
	}

	return gorm.G[model.ApprovalRoom](ar.db).Create(ctx, newRoom)
}

func (*ApprovalRoomRepository) CreateWithTx(
	title string,
	filepaths []string,
	dueAt time.Time,
	submitterId string,
	approvalRoomId *string,
) func(ctx context.Context, tx *gorm.DB) error {
	roomID := uuid.New().String()
	*approvalRoomId = roomID

	newRoom := &model.ApprovalRoom{
		ID:          roomID,
		Title:       title,
		Filepaths:   strings.Join(filepaths, ";"),
		DueAt:       dueAt,
		SubmitterId: submitterId,
	}

	return func(ctx context.Context, tx *gorm.DB) error {
		ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
		defer cancel()

		return tx.Create(newRoom).Error
	}
}
