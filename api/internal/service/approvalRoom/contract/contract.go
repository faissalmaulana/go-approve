package contract

import (
	"time"

	"github.com/faissalmaulana/go-approve/internal/model"
)

type CreateApprovalRoom struct {
	Title         string
	DueAt         time.Time
	SubmitterId   string
	InviteeIds    []string
	FileMetadatas map[string]model.FileMetadata
}

type ApprovalDocument struct {
	Link            string `json:"link"`
	DisplayFileName string `json:"display_file_name"`
	Size            int    `json:"size"`
}

type ApprovalApprover struct {
	Handle   string `json:"handle"`
	Name     string `json:"name"`
	Decision string `json:"decision"`
}

type ApprovalAggregates struct {
	FileUploaded     int `json:"file_uploaded"`
	ApprovalProgress int `json:"approval_progress"`
}

type GetApprovalRoomByID struct {
	Title           string             `json:"title"`
	CreatedAt       time.Time          `json:"created_at"`
	DueAt           time.Time          `json:"due_at"`
	SubmitterHandle string             `json:"submitter_handle"`
	Documents       []ApprovalDocument `json:"documents"`
	Approvers       []ApprovalApprover `json:"approvers"`
	Aggregates      ApprovalAggregates `json:"aggregates"`
}

type ApprovalRoomRequest struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	DueAt     time.Time `json:"due_at"`
	CreatedAt time.Time `json:"created_at"`
}
