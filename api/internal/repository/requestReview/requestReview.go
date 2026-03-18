package requestreview

import (
	"context"

	"github.com/faissalmaulana/go-approve/internal/constant"
	"github.com/faissalmaulana/go-approve/internal/model"
	"github.com/faissalmaulana/go-approve/internal/utils"
	"gorm.io/gorm"
)

type RequestReviewStorage interface {
	CreateBatchWithTx(
		inviteeIds []string,
		requesterId string,
		approvalRoomId *string,
	) func(ctx context.Context, tx *gorm.DB) error
	UpdateWithTx(id string, status utils.Status) func(ctx context.Context, tx *gorm.DB) error
	GetById(ctx context.Context, id string) (model.ReviewRequest, error)
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

func (r *RequestReviewRepository) GetById(ctx context.Context, id string) (model.ReviewRequest, error) {
	return gorm.G[model.ReviewRequest](r.db).Where("id = ?", id).First(ctx)
}

func (r *RequestReviewRepository) UpdateWithTx(id string, status utils.Status) func(ctx context.Context, tx *gorm.DB) error {

	return func(ctx context.Context, tx *gorm.DB) error {
		ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
		defer cancel()

		_, err := gorm.G[model.ReviewRequest](tx).Where("id = ?", id).Update(ctx, "status", status.String())
		return err
	}
}
