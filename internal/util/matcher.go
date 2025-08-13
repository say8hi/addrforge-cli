package util

import "strings"

// Match performs case-insensitive prefix matching
func Match(addr string, prefix string) bool {
	return strings.HasPrefix(strings.ToLower(addr), strings.ToLower(prefix))
}

// MatchCaseSensitive performs case-sensitive prefix matching
func MatchCaseSensitive(addr string, prefix string) bool {
	return strings.HasPrefix(addr, prefix)
}
