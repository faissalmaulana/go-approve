package approvalroom

import (
	"context"

	approvalroomrepository "github.com/faissalmaulana/go-approve/internal/repository/approvalRoom"
	requestreviewrepository "github.com/faissalmaulana/go-approve/internal/repository/requestReview"
	"github.com/faissalmaulana/go-approve/internal/repository/transactions"
	"github.com/faissalmaulana/go-approve/internal/service/approvalRoom/contract"
)

type ApprovalRoomService struct {
	dbTransaction        transactions.DatabaseTransaction
	approvalRoomStorage  approvalroomrepository.ApprovalRoomStorage
	requestReviewStorage requestreviewrepository.RequestReviewStorage
}

func New(
	ars approvalroomrepository.ApprovalRoomStorage,
	rrs requestreviewrepository.RequestReviewStorage,
	dbtx transactions.DatabaseTransaction,
) *ApprovalRoomService {
	return &ApprovalRoomService{
		approvalRoomStorage:  ars,
		requestReviewStorage: rrs,
		dbTransaction:        dbtx,
	}
}

func (a *ApprovalRoomService) Create(ctx context.Context, i *contract.CreateApprovalRoom) error {
	var approvalRoomId string

	createApprovalRoom := a.approvalRoomStorage.CreateWithTx(
		i.Title,
		i.Filepaths,
		i.DueAt,
		i.SubmitterId,
		&approvalRoomId,
	)

	createReviewRequests := a.requestReviewStorage.CreateBatchWithTx(
		i.InviteeIds,
		i.SubmitterId,
		&approvalRoomId,
	)

	return a.dbTransaction.RunTransactions(ctx, createApprovalRoom, createReviewRequests)
}
