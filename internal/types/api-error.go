package types

import "fmt"

// APIError represents a structured error message for API responses.
// It contains the error code, message, the operation during which the error occurred,
// and any embedded error details for deeper debugging.
type APIError struct {
	Code          int    `json:"code"`          // Code is the HTTP status code associated with the error.
	Message       string `json:"message"`       // Message is a human-readable description of the error.
	Operation     string `json:"operation"`     // Operation specifies the HTTP method and path where the error occurred.
	EmbeddedError string `json:"embeddedError"` // EmbeddedError provides additional error information, if available.
}

// Error implements the error interface for APIError, returning an empty string.
// The error message can be customized to include more detail or structured formatting as needed.
//
// Returns:
// - A string representing the error message (currently empty).
func (err *APIError) Error() string {
	return ""
}

// FormatOperation formats the HTTP method and path into a single string representing the operation.
// This function is useful for creating the Operation field in APIError.
//
// Parameters:
// - method: The HTTP method (e.g., "GET", "POST").
// - path: The HTTP path where the operation is performed.
//
// Returns:
// - A string combining the method and path in the format "METHOD /path".
func FormatOperation(method, path string) string {
	return fmt.Sprintf("%s %s", method, path)
}
