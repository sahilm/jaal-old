package jaal

import (
	"crypto/sha256"
	"encoding/hex"
)

func ToSHA256(s string) (string, error) {
	h := sha256.New()
	_, err := h.Write([]byte(s))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
