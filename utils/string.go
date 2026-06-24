package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func CalculateURIHash(uri string) string {
	h := sha256.New()
	_, _ = h.Write([]byte(uri))
	hash := h.Sum(nil)
	return hex.EncodeToString(hash)
}
