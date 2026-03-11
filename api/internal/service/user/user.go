package user

import (
	"context"

	"github.com/faissalmaulana/go-approve/internal/model/public"
	"github.com/faissalmaulana/go-approve/internal/repository/user"
	"github.com/faissalmaulana/go-approve/internal/service"
	"gorm.io/gorm"
)

type User struct {
	userStorage user.UserStorage
}

func New(u user.UserStorage) *User {
	return &User{
		userStorage: u,
	}
}

func (u *User) GetCurrentUser(ctx context.Context, id string) (*public.UserPublic, error) {
	if id == "" {
		return nil, service.ErrSubIsEmpty
	}

	user, err := u.userStorage.FindByID(ctx, id)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, service.ErrUserNotFound
		default:
			return nil, service.ErrInternal
		}
	}

	return user, nil
}
