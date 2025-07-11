package subsidy

import "errors"

var (
	ErrTransactionFailed  = errors.New("blockchain transaction failed")
	ErrInvalidInput       = errors.New("invalid input parameters")
	ErrNotFound           = errors.New("resource not found")
	ErrTimeout            = errors.New("operation timed out")
	ErrDistributionFailed = errors.New("subsidy distribution failed")
	ErrInvalidEpochState  = errors.New("epoch is not in valid state for operation")
)
