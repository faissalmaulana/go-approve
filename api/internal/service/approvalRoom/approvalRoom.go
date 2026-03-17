package approvalroom

import (
	"context"
	"time"

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
	return &contract.GetApprovalRoomByID{
		Title:           "FIFA World Cup 2026 Sponsorship Agreement",
		CreatedAt:       time.Now(),
		DueAt:           time.Now().Add(time.Hour * 72),
		SubmitterHandle: "John Doe",
		Documents: []contract.ApprovalDocument{
			{
				Link:            "localhost:8080/example.com",
				DisplayFileName: "Sponsorship Agreement v3.pdf",
				Size:            2048,
			},
			{
				Link:            "localhost:8080/example.com",
				DisplayFileName: "Budget Breakdown Q1.xlsx",
				Size:            512,
			},
			{
				Link:            "localhost:8080/example.com",
				DisplayFileName: "Legal Terms Final.pdf",
				Size:            1024,
			},
		},
		Approvers: []contract.ApprovalApprover{
			{
				Handle:   "@sarah.connor",
				Name:     "Sarah Connor",
				Decision: "approved",
			},
			{
				Handle:   "@michael.scott",
				Name:     "Michael Scott",
				Decision: "pending",
			},
			{
				Handle:   "@tony.stark",
				Name:     "Tony Stark",
				Decision: "rejected",
			},
		},
		Aggregates: contract.ApprovalAggregates{
			FileUploaded:     3,
			ApprovalProgress: 33,
		},
	}, nil
}
