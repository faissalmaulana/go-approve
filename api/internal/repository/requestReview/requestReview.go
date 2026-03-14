package requestreview

import (
	"context"

	"github.com/faissalmaulana/go-approve/internal/constant"
	"github.com/faissalmaulana/go-approve/internal/model"
	"gorm.io/gorm"
)

type RequestReviewStorage interface {
	CreateBatchWithTx(
		inviteeIds []string,
		requesterId string,
		approvalRoomId *string,
	) func(ctx context.Context, tx *gorm.DB) error
}

type RequestReviewRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) RequestReviewStorage {
	return &RequestReviewRepository{
		db: db,
	}
}

func (r *RequestReviewRepository) CreateBatchWithTx(
	inviteeIds []string,
	requesterId string,
	approvalRoomId *string,
) func(ctx context.Context, tx *gorm.DB) error {

	reviewRequests := make([]*model.ReviewRequest, 0, len(inviteeIds))
	for _, inviteeId := range inviteeIds {
		reviewRequests = append(reviewRequests, &model.ReviewRequest{
			Status:         "pending",
			ApprovalRoomId: *approvalRoomId,
			InviteeId:      inviteeId,
			RequesterId:    requesterId,
		})
	}

	return func(ctx context.Context, tx *gorm.DB) error {
		ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
		defer cancel()

		return tx.Create(reviewRequests).Error
	}
}
