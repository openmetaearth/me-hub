package utils

import "strings"

const BridgeTokenPrefix = "u"

func GetDenom(denom string) string {
	return BridgeTokenPrefix + strings.ToLower(strings.TrimSpace(denom))
}
