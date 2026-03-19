package dto

type GetRequestReviewDTO struct {
	As     string `query:"as" validate:"omitempty,oneof=sent received"`
	Status string `query:"status" validate:"omitempty,oneof=pending accepted rejected"`
	Limit  int    `query:"limit" validate:"omitempty,gte=0,lte=100"`
	Offset int    `query:"offset" validate:"omitempty,gte=0"`
}

