package approvalroomapprover

import (
	"context"

	"github.com/faissalmaulana/go-approve/internal/constant"
	"github.com/faissalmaulana/go-approve/internal/model"
	"gorm.io/gorm"
)

type ApprovalRoomApproverStorage interface {
	CreateWithTx(approvalId, approvalRoomId string) func(ctx context.Context, tx *gorm.DB) error
}

type ApprovalRoomApproverRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) ApprovalRoomApproverStorage {
	return &ApprovalRoomApproverRepository{
		db: db,
	}
}

func (a *ApprovalRoomApproverRepository) CreateWithTx(approvalId, approvalRoomId string) func(ctx context.Context, tx *gorm.DB) error {
	newApprovalRoomApprover := &model.ApprovalRoomApprover{
		ApprovalId:     approvalId,
		ApprovalRoomId: approvalRoomId,
		Decision:       "pending",
	}

	return func(ctx context.Context, tx *gorm.DB) error {
		ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
		defer cancel()

		return tx.Create(newApprovalRoomApprover).Error
	}
}
