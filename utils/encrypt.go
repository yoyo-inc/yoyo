package utils

import (
	"crypto/sha1"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

// Encrypt encrypts plain text by md5sum
func Encrypt(plainText string) string {
	var val = []byte(plainText)
	return fmt.Sprintf("%x", pbkdf2.Key(val, []byte("VglXSIH9WGEwkdB3"), 4096, 32, sha1.New))
}
