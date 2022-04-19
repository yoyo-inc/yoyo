package utils

import (
	"crypto/md5"
	"fmt"
)

// Encrypt encrypts plain text by md5sum
func Encrypt(plainText string) string {
	var val = []byte(plainText)
	return fmt.Sprintf("%x", md5.Sum(val))
}
