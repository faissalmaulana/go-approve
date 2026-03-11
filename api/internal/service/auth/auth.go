package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/faissalmaulana/go-approve/internal/repository/user"
	"github.com/faissalmaulana/go-approve/internal/service"
	"github.com/faissalmaulana/go-approve/internal/service/jwtfx"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Auth struct {
	userStore user.UserStorage
	jwt       jwtfx.TokenGenerator
}

func New(u user.UserStorage, j jwtfx.TokenGenerator) *Auth {
	return &Auth{
		userStore: u,
		jwt:       j,
	}
}

func (a *Auth) Register(ctx context.Context, email, name, password, handler string) (string, error) {
	var emptyFields []string
	if email == "" {
		emptyFields = append(emptyFields, "email")
	}
	if name == "" {
		emptyFields = append(emptyFields, "name")
	}
	if password == "" {
		emptyFields = append(emptyFields, "password")
	}
	if handler == "" {
		emptyFields = append(emptyFields, "handler")
	}

	if len(emptyFields) > 0 {
		return "", fmt.Errorf("%w: %s", service.ErrInvalidPayload, strings.Join(emptyFields, ", "))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	newUserId, err := a.userStore.Create(
		ctx,
		email,
		name,
		string(hashedPassword),
		handler)

	if err != nil {
		switch err {
		case gorm.ErrDuplicatedKey:
			return "", service.ErrDuplicatedUser
		default:
			return "", service.ErrInternal
		}
	}

	accessToken, err := a.jwt.GenerateAccessToken(newUserId)
	if err != nil {
		return "", service.ErrInternal
	}

	return accessToken, nil
}

func (a *Auth) Login(ctx context.Context, email, password string) (string, error) {
	var emptyFields []string
	if email == "" {
		emptyFields = append(emptyFields, "email")
	}

	if password == "" {
		emptyFields = append(emptyFields, "password")
	}

	if len(emptyFields) > 0 {
		return "", fmt.Errorf("%w: %s", service.ErrInvalidPayload, strings.Join(emptyFields, ", "))
	}

	user, err := a.userStore.FindByEmail(ctx, email)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return "", service.ErrUserNotFound
		default:
			return "", service.ErrInternal
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", service.ErrPasswordNotMatched
	}

	accessToken, err := a.jwt.GenerateAccessToken(user.ID)
	if err != nil {
		return "", service.ErrInternal
	}

	return accessToken, nil
}
