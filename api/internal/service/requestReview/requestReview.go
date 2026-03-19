package requestreview

import (
	"context"

	approvalroomapprover "github.com/faissalmaulana/go-approve/internal/repository/approvalRoomApprover"
	requestreview "github.com/faissalmaulana/go-approve/internal/repository/requestReview"
	"github.com/faissalmaulana/go-approve/internal/repository/transactions"
	"github.com/faissalmaulana/go-approve/internal/model"
	"github.com/faissalmaulana/go-approve/internal/utils"
)

type RequestReviewService struct {
	RequestReviewStorage    requestreview.RequestReviewStorage
	ApprovalApproverStorage approvalroomapprover.ApprovalRoomApproverStorage
	dbTransaction           transactions.DatabaseTransaction
}

func New(rrs requestreview.RequestReviewStorage, aas approvalroomapprover.ApprovalRoomApproverStorage, dbtx transactions.DatabaseTransaction) *RequestReviewService {
	return &RequestReviewService{
		RequestReviewStorage:    rrs,
		ApprovalApproverStorage: aas,
		dbTransaction:           dbtx,
	}
}

func (r *RequestReviewService) ConfirmRequestRevieWithInsertApprover(
	ctx context.Context,
	requestReviewId,
	approvalId string,
	status utils.Status,
) error {

	requestReview, err := r.RequestReviewStorage.GetById(ctx, requestReviewId)
	if err != nil {
		return err
	}

	requestReviewUpdateStatusTx := r.RequestReviewStorage.UpdateWithTx(requestReview.ID, status)

	createApprovalApproverTx := r.ApprovalApproverStorage.CreateWithTx(approvalId, requestReview.ApprovalRoomId)

	return r.dbTransaction.RunTransactions(ctx, requestReviewUpdateStatusTx, createApprovalApproverTx)
}

func (r *RequestReviewService) GetReceivedInvitations(
	ctx context.Context,
	inviteeId string,
	status *utils.Status,
	limit int,
	offset int,
) ([]model.ReviewRequest, error) {
	return r.RequestReviewStorage.GetReceivedInvitations(ctx, inviteeId, status, limit, offset)
}
