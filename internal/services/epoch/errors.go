package epoch

import "errors"

// Predefined error types for different failure scenarios
var (
	ErrTransactionFailed = errors.New("blockchain transaction failed")
	ErrInvalidInput      = errors.New("invalid input parameters")
	ErrNotFound          = errors.New("resource not found")
	ErrTimeout           = errors.New("operation timed out")
)
