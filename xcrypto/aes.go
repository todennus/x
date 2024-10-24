package xcrypto

import (
	"crypto/sha256"

	"golang.org/x/crypto/pbkdf2"
)

func GenerateAESKeyFromPassword(password string, size int) []byte {
	return pbkdf2.Key([]byte(password), nil, 10000, size, sha256.New)
}
