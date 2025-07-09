package handlers

import (
	"errors"
	"net/http"

	"github.com/andrey/epoch-server/internal/services/epoch"
	"github.com/andrey/epoch-server/internal/services/merkle"
	"github.com/andrey/epoch-server/internal/services/subsidy"
	"github.com/go-pkgz/lgr"
	"github.com/go-pkgz/rest"
)

// ErrorResponse represents the structure of error responses
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}

// writeErrorResponse writes a structured error response based on the error type
func writeErrorResponse(w http.ResponseWriter, r *http.Request, logger lgr.L, err error, message string) {
	// Determine appropriate HTTP status code based on error type
	var statusCode int
	if isTransactionFailedError(err) {
		statusCode = http.StatusBadGateway
	} else if isInvalidInputError(err) {
		statusCode = http.StatusBadRequest
	} else if isNotFoundError(err) {
		statusCode = http.StatusNotFound
	} else if isTimeoutError(err) {
		statusCode = http.StatusRequestTimeout
	} else {
		// Default to internal server error
		statusCode = http.StatusInternalServerError
	}

	rest.SendErrorJSON(w, r, logger, statusCode, err, message)
}

// Helper functions to check error types across all services
func isTransactionFailedError(err error) bool {
	return errors.Is(err, epoch.ErrTransactionFailed) ||
		   errors.Is(err, subsidy.ErrTransactionFailed)
}

func isInvalidInputError(err error) bool {
	return errors.Is(err, epoch.ErrInvalidInput) ||
		   errors.Is(err, subsidy.ErrInvalidInput) ||
		   errors.Is(err, merkle.ErrInvalidInput)
}

func isNotFoundError(err error) bool {
	return errors.Is(err, epoch.ErrNotFound) ||
		   errors.Is(err, subsidy.ErrNotFound) ||
		   errors.Is(err, merkle.ErrNotFound)
}

func isTimeoutError(err error) bool {
	return errors.Is(err, epoch.ErrTimeout) ||
		   errors.Is(err, subsidy.ErrTimeout)
}