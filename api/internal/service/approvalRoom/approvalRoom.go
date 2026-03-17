package approvalroom

import (
	"context"

	approvalroomrepository "github.com/faissalmaulana/go-approve/internal/repository/approvalRoom"
	filemetadatarepository "github.com/faissalmaulana/go-approve/internal/repository/fileMetadata"
	requestreviewrepository "github.com/faissalmaulana/go-approve/internal/repository/requestReview"
	"github.com/faissalmaulana/go-approve/internal/repository/transactions"
	"github.com/faissalmaulana/go-approve/internal/service/approvalRoom/contract"
)

type ApprovalRoomService struct {
	dbTransaction        transactions.DatabaseTransaction
	approvalRoomStorage  approvalroomrepository.ApprovalRoomStorage
	requestReviewStorage requestreviewrepository.RequestReviewStorage
	fileMetadataStorage  filemetadatarepository.FileMetadataStorage
}

func New(
	ars approvalroomrepository.ApprovalRoomStorage,
	rrs requestreviewrepository.RequestReviewStorage,
	fms filemetadatarepository.FileMetadataStorage,
	dbtx transactions.DatabaseTransaction,
) *ApprovalRoomService {
	return &ApprovalRoomService{
		approvalRoomStorage:  ars,
		requestReviewStorage: rrs,
		fileMetadataStorage:  fms,
		dbTransaction:        dbtx,
	}
}

func (a *ApprovalRoomService) Create(ctx context.Context, i *contract.CreateApprovalRoom) (string, error) {
	var approvalRoomId string

	createApprovalRoom := a.approvalRoomStorage.CreateWithTx(
		i.Title,
		i.DueAt,
		i.SubmitterId,
		&approvalRoomId,
	)

	createReviewRequests := a.requestReviewStorage.CreateBatchWithTx(
		i.InviteeIds,
		i.SubmitterId,
		&approvalRoomId,
	)

	createFileMetadatas := a.fileMetadataStorage.CreateBatchWithTx(
		i.FileMetadatas,
		approvalRoomId,
	)

	if err := a.dbTransaction.RunTransactions(ctx, createApprovalRoom, createReviewRequests, createFileMetadatas); err != nil {
		return "", err
	}

	return approvalRoomId, nil
}
