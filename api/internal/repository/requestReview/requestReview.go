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
	Update(ctx context.Context, id string, status utils.Status) error
	GetById(ctx context.Context, id string) (model.ReviewRequest, error)
	GetReceivedInvitations(
		ctx context.Context,
		inviteeId string,
		status *utils.Status,
		limit int,
		offset int,
	) ([]model.ReviewRequest, error)
	GetSentInvitations(
		ctx context.Context,
		requesterId string,
		status *utils.Status,
		limit int,
		offset int,
	) ([]model.ReviewRequest, error)
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

func (r *RequestReviewRepository) Update(ctx context.Context, id string, status utils.Status) error {
	ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
	defer cancel()

	_, err := gorm.G[model.ReviewRequest](r.db).Where("id = ?", id).Update(ctx, "status", status.String())
	return err
}

func (r *RequestReviewRepository) UpdateWithTx(id string, status utils.Status) func(ctx context.Context, tx *gorm.DB) error {

	return func(ctx context.Context, tx *gorm.DB) error {
		ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
		defer cancel()

		_, err := gorm.G[model.ReviewRequest](tx).Where("id = ?", id).Update(ctx, "status", status.String())
		return err
	}
}

func (r *RequestReviewRepository) GetReceivedInvitations(
	ctx context.Context,
	inviteeId string,
	status *utils.Status,
	limit int,
	offset int,
) ([]model.ReviewRequest, error) {
	ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
	defer cancel()

	q := gorm.G[model.ReviewRequest](r.db).
		Where("invitee_id = ?", inviteeId).
		Preload("Requester", nil).
		Preload("ApprovalRoom", nil)

	if status != nil {
		q = q.Where("status = ?", status.String())
	}

	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}

	return q.Order("created_at desc").Find(ctx)
}

func (r *RequestReviewRepository) GetSentInvitations(
	ctx context.Context,
	requesterId string,
	status *utils.Status,
	limit int,
	offset int,
) ([]model.ReviewRequest, error) {
	ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
	defer cancel()

	q := gorm.G[model.ReviewRequest](r.db).
		Where("requester_id = ?", requesterId).
		Preload("Invitee", nil).
		Preload("ApprovalRoom", nil)

	if status != nil {
		q = q.Where("status = ?", status.String())
	}

	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}

	return q.Order("created_at desc").Find(ctx)
}
