package jwtfx

import (
	"time"

	"github.com/faissalmaulana/go-approve/internal/service"
	gojwt "github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"
)

type TokenGenerator interface {
	GenerateAccessToken(sub string) (string, error)
}

type JwtFx struct {
	config *Config
}

var Module = fx.Module("jwt", fx.Provide(New), fx.Provide(fx.Private, NewConfig))

func New(c *Config) TokenGenerator {
	return &JwtFx{
		config: c,
	}
}

func (j *JwtFx) GenerateAccessToken(sub string) (string, error) {
	if sub == "" {
		return "", service.ErrSubIsEmpty
	}

	claims := gojwt.RegisteredClaims{
		Subject:   sub,
		Issuer:    j.config.Issuer,
		ExpiresAt: gojwt.NewNumericDate(time.Now().Add(j.config.Expire)),
	}

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.config.SecretKey))

}
