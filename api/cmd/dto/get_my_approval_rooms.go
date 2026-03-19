package dto

type GetMyApprovalRoomsDTO struct {
	Sort  string `query:"sort" validate:"omitempty,oneof=due_at created_at"`
	Order string `query:"order" validate:"omitempty,oneof=asc desc"`

	Limit  int `query:"limit" validate:"omitempty,gte=0,lte=100"`
	Offset int `query:"offset" validate:"omitempty,gte=0"`
}

