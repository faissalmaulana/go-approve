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

func (a *ApprovalRoomService) GetApprovalRoomById(id string) (*contract.GetApprovalRoomByID, error) {
	ctx := context.Background()

	detail, err := a.approvalRoomStorage.GetApprovalRoomByID(ctx, id)
	if err != nil {
		return nil, err
	}

	counts, err := a.approvalRoomStorage.GetApprovalRoomCountsByID(ctx, id)
	if err != nil {
		return nil, err
	}

	documents := make([]contract.ApprovalDocument, 0, len(detail.Files))
	for _, f := range detail.Files {
		documents = append(documents, contract.ApprovalDocument{
			Link:            f.Link,
			DisplayFileName: f.Filename,
			Size:            f.Size,
		})
	}

	approvers := make([]contract.ApprovalApprover, 0, len(detail.Approvers))
	for _, ap := range detail.Approvers {
		approvers = append(approvers, contract.ApprovalApprover{
			Handle:   ap.Handle,
			Name:     ap.Name,
			Decision: ap.Decision,
		})
	}

	approvalProgress := 0
	if counts.TotalApprover > 0 {
		approvalProgress = int(counts.ApprovedCount * 100 / counts.TotalApprover)
	}

	return &contract.GetApprovalRoomByID{
		Title:           detail.Room.Title,
		CreatedAt:       detail.Room.CreatedAt,
		DueAt:           detail.Room.DueAt,
		SubmitterHandle: detail.Submitter.Handler,
		Documents:       documents,
		Approvers:       approvers,
		Aggregates: contract.ApprovalAggregates{
			FileUploaded:     int(counts.FileUploaded),
			ApprovalProgress: approvalProgress,
		},
	}, nil
}
