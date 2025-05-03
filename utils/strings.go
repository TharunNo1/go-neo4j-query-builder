package utils

import "strings"

// JoinStrings joins a slice of strings with a separator
func JoinStrings(strs []string, sep string) string {
	return strings.Join(strs, sep)
}
