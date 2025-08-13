package util

import "strings"

func Match(addr string, prefix string) bool {
	return strings.HasPrefix(strings.ToLower(addr), strings.ToLower(prefix))
}

func MatchCaseSensetive(addr string, prefix string) bool {
	return strings.HasPrefix(addr, prefix)
}
