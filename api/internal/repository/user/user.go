package user

import (
	"context"

	"github.com/faissalmaulana/go-approve/internal/constant"
	"github.com/faissalmaulana/go-approve/internal/model"
	"gorm.io/gorm"
)

type UserStorage interface {
	Create(ctx context.Context, email, name, password, handler string) (string, error)
}

type UserRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) UserStorage {
	return &UserRepository{
		db,
	}
}

func (r *UserRepository) Create(ctx context.Context, email, name, password, handler string) (string, error) {
	u := &model.User{
		Email:    email,
		Name:     name,
		Handler:  handler,
		Password: password,
	}

	ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
	defer cancel()

	if err := gorm.G[model.User](r.db).Create(ctx, u); err != nil {
		return "", err
	}

	return u.ID, nil
}
