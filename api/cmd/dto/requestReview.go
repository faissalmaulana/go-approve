package dto

type ConfirmReviewRequestDTO struct {
	Status string `json:"status" validate:"required,oneof=accepted rejected"`
}
