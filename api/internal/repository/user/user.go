package user

import (
	"context"

	"github.com/faissalmaulana/go-approve/internal/constant"
	"github.com/faissalmaulana/go-approve/internal/model"
	"github.com/faissalmaulana/go-approve/internal/model/public"
	"gorm.io/gorm"
)

type UserStorage interface {
	Create(ctx context.Context, email, name, password, handler string) (string, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id string) (*public.UserPublic, error)
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

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
	defer cancel()

	user, err := gorm.G[model.User](r.db).Where("email = ?", email).First(ctx)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*public.UserPublic, error) {
	ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
	defer cancel()

	var user = new(public.UserPublic)

	err := gorm.G[model.User](r.db).Select("id", "name", "handler", "email").Where("id = ?", id).Limit(1).Scan(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
