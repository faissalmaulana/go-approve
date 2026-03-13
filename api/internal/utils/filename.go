package utils

import (
	"crypto/rand"
	"encoding/hex"
	"path/filepath"
)

func GenerateRandomFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)

	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
	randomString := hex.EncodeToString(randomBytes)

	return randomString + ext
}
