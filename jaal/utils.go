package jaal

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"time"
)

func ToSHA256(s string) (string, error) {
	h := sha256.New()
	_, err := h.Write([]byte(s))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func EnrichEvent(event *Event) {
	now := time.Now()
	event.SourceHostName = lookupAddr(event.Source)
	event.UnixTime = now.Unix()
	event.Timestamp = now.UTC().Format(time.RFC3339)
}

func lookupAddr(address string) string {
	hosts, err := net.LookupAddr(address)
	if err != nil {
		return "" // Don't care on err, just return nothing
	}
	return hosts[0]
}
