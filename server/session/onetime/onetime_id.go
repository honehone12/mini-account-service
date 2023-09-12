package onetime

import (
	"crypto/rand"
	"encoding/base64"
)

const (
	OnetimeIdBytesLength = 32
)

func New() (string, error) {
	buff := make([]byte, OnetimeIdBytesLength)
	_, err := rand.Read(buff)
	if err != nil {
		return "", err
	}
	id := base64.StdEncoding.EncodeToString(buff)
	return id, nil
}
