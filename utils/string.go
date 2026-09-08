package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

func CalculateURIHash(uri string) string {
	h := sha256.New()
	_, _ = io.WriteString(h, uri)
	return hex.EncodeToString(h.Sum(nil))
}
