package merkle

import "errors"

// Predefined error types for merkle proof operations
var (
	ErrInvalidInput    = errors.New("invalid input parameters")
	ErrNotFound        = errors.New("resource not found")
	ErrProofGeneration = errors.New("merkle proof generation failed")
	ErrInvalidProof    = errors.New("invalid merkle proof")
)
