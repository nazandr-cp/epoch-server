package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/andrey/epoch-server/internal/services/epoch"
	"github.com/andrey/epoch-server/internal/services/merkle"
	"github.com/andrey/epoch-server/internal/services/subsidy"
)

// ErrorResponse represents the structure of error responses
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}

// writeErrorResponse writes a structured error response based on the error type
func writeErrorResponse(w http.ResponseWriter, err error, message string) {
	w.Header().Set("Content-Type", "application/json")
	
	var errResponse ErrorResponse
	errResponse.Error = message
	errResponse.Details = err.Error()

	// Determine appropriate HTTP status code based on error type
	if isTransactionFailedError(err) {
		errResponse.Code = http.StatusBadGateway
		w.WriteHeader(http.StatusBadGateway)
	} else if isInvalidInputError(err) {
		errResponse.Code = http.StatusBadRequest
		w.WriteHeader(http.StatusBadRequest)
	} else if isNotFoundError(err) {
		errResponse.Code = http.StatusNotFound
		w.WriteHeader(http.StatusNotFound)
	} else if isTimeoutError(err) {
		errResponse.Code = http.StatusRequestTimeout
		w.WriteHeader(http.StatusRequestTimeout)
	} else {
		// Default to internal server error
		errResponse.Code = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(errResponse)
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