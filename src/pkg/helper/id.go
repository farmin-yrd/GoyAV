package helper

import (
	"crypto/md5"
	"encoding/base64"
	"io"
)

func MakeID(b io.Reader) (string, error) {
	hash := md5.New()
	_, err := io.Copy(hash, b)
	if err != nil {
		return "", err
	}
	id := base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
	return id, nil
}

func IsValidID(id string) bool {
	return id != ""
}
