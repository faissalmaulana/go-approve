package dto

import "time"

type CreateApprovalRoomDTO struct {
	Title     string    `json:"title" validate:"required,min=6,max=255"`
	DueAt     time.Time `json:"due_at" validate:"required"`
	Approvers []string  `json:"approvers" validate:"required,gt=0,dive,required"`
}
