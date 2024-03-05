package helper

import (
	"crypto/sha256"
	"fmt"
	"io"
)

func MakeHash(i io.Reader) (string, error) {
	hash := sha256.New()
	_, err := io.Copy(hash, i)
	if err != nil {
		return "", err
	}
	// Calculate the SHA-256 hash of the content.
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func IsValidHash(h string) bool {
	return h != ""
}
