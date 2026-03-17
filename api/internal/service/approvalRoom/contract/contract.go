package contract

import (
	"time"

	"github.com/faissalmaulana/go-approve/internal/model"
)

type CreateApprovalRoom struct {
	Title         string
	DueAt         time.Time
	Filepaths     []string
	SubmitterId   string
	InviteeIds    []string
	FileMetadatas map[string]model.FileMetadata
}
