package approvalroom

import (
	"context"

	approvalroomrepository "github.com/faissalmaulana/go-approve/internal/repository/approvalRoom"
	"github.com/faissalmaulana/go-approve/internal/service/approvalRoom/contract"
)

type ApprovalRoomService struct {
	approvalRoomStorage approvalroomrepository.ApprovalRoomStorage
}

func New(ars approvalroomrepository.ApprovalRoomStorage) *ApprovalRoomService {
	return &ApprovalRoomService{
		approvalRoomStorage: ars,
	}
}

func (a *ApprovalRoomService) Create(ctx context.Context, i *contract.CreateApprovalRoom) error {

	// TODO: filter error message
	return a.approvalRoomStorage.Create(
		ctx,
		i.Title,
		i.Filepaths,
		i.DueAt,
		i.SubmitterId,
	)
}
