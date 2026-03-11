package jwtfx

import (
	"testing"
	"time"

	"github.com/faissalmaulana/go-approve/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAccessToken(t *testing.T) {
	j := JwtFx{
		config: &Config{
			SecretKey: "super_secret_key",
			Issuer:    "unit_test",
			Expire:    time.Minute * 1,
		},
	}

	t.Run("success generate access token and return the generated access token", func(t *testing.T) {
		subject := "test"
		result, err := j.GenerateAccessToken(subject)
		assert.NoError(t, err)

		assert.NotEmpty(t, result)
	})

	t.Run("error generate access token because subject is empty string", func(t *testing.T) {
		result, err := j.GenerateAccessToken("")

		assert.ErrorIs(t, err, service.ErrSubIsEmpty)
		assert.Emptyf(t, result, `GenerateAccessToken("") should return empty string`)
	})

}
