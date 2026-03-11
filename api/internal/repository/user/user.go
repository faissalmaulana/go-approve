package user

import (
	"context"

	"github.com/faissalmaulana/go-approve/internal/constant"
	"github.com/faissalmaulana/go-approve/internal/model"
	"gorm.io/gorm"
)

type UserStorage interface {
	Create(ctx context.Context, email, name, password, handler string) error
}

type UserRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) UserStorage {
	return &UserRepository{
		db,
	}
}

func (r *UserRepository) Create(ctx context.Context, email, name, password, handler string) error {
	u := &model.User{
		Email:    email,
		Name:     name,
		Handler:  handler,
		Password: password,
	}

	ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
	defer cancel()

	return gorm.G[model.User](r.db).Create(ctx, u)
}
