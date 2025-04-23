package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func CalculateUriHash(uri string) string {
	h := sha256.New()
	h.Write([]byte(uri))
	hash := h.Sum(nil)
	return hex.EncodeToString(hash)
}
