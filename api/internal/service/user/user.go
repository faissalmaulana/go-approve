package user

import (
	"context"

	"github.com/faissalmaulana/go-approve/internal/model/public"
	"github.com/faissalmaulana/go-approve/internal/repository/blocklisttoken"
	"github.com/faissalmaulana/go-approve/internal/repository/user"
	"github.com/faissalmaulana/go-approve/internal/service"
	"gorm.io/gorm"
)

type User struct {
	userStorage   user.UserStorage
	blocklistRepo blocklisttoken.BlocklistStorage
}

func New(u user.UserStorage, b blocklisttoken.BlocklistStorage) *User {
	return &User{
		userStorage:   u,
		blocklistRepo: b,
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

func (u *User) Logout(ctx context.Context, token string) error {
	if token == "" {
		return service.ErrInvalidPayload
	}

	if err := u.blocklistRepo.Create(ctx, token); err != nil {
		return service.ErrInternal
	}

	return nil
}

func (u *User) GetUserIdsOnly(ctx context.Context, ids []string) ([]string, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	users, err := u.userStorage.FindByIDs(ctx, ids)
	if err != nil {
		return nil, service.ErrInternal
	}

	result := make([]string, len(users))
	for i, user := range users {
		result[i] = user.ID
	}

	return result, nil
}

func (u *User) SearchUsers(ctx context.Context, handle string) (*[]public.UserPublic, error) {
	users, err := u.userStorage.SearchUsersByHandle(ctx, handle)
	if err != nil {
		return nil, service.ErrInternal
	}

	return users, nil
}
