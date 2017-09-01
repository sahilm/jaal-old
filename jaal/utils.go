package jaal

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
)

func ToSHA256(b []byte) (string, error) {
	h := sha256.New()
	_, err := h.Write(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func LookupAddr(address string) string {
	hosts, err := net.LookupAddr(address)
	if err != nil {
		return "" // Don't care on err, just return nothing
	}
	return hosts[0]
}
