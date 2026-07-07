package utils

import (
	"fmt"
)

func MustStrToByte32(str string) [32]byte {
	byte32, err := StrToByte32(str)
	if err != nil {
		panic(err)
	}
	return byte32
}

func Byte32ToString(bytes [32]byte) string {
	for i := len(bytes) - 1; i >= 0; i-- {
		if bytes[i] != 0 {
			return string(bytes[:i+1])
		}
	}
	return ""
}

func StrToByte32(s string) ([32]byte, error) {
	var out [32]byte
	if len([]byte(s)) > 32 {
		return out, fmt.Errorf("string too long")
	}
	copy(out[:], s)
	return out, nil
}
