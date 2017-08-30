package jaal

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"os"
)

func ToSHA256(s string) (string, error) {
	h := sha256.New()
	_, err := h.Write([]byte(s))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func LookupAddr(address string) string {
	ip, _, err := net.SplitHostPort(address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to lookup %v. Error: %v", address, err)
		return "" // Don't care on err, just return nothing
	}
	hosts, err := net.LookupAddr(ip)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to lookup %v. Error: %v", address, err)
		return "" // Don't care on err, just return nothing
	}
	return hosts[0]
}
