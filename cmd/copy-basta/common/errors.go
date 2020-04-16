package common

import "fmt"

// The NewFlagValidationError build an error that mimics cobra flag error formatting
func NewFlagValidationError(flag string, reason string) error {
	return fmt.Errorf(`"--%s" %s`, flag, reason)
}
