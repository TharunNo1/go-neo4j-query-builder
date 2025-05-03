package utils

import "fmt"

// ValidateString checks if the given string is non-empty
func ValidateString(s string) error {
	if s == "" {
		return fmt.Errorf("string cannot be empty")
	}
	return nil
}
