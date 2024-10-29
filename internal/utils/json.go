package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ParseJSON decodes a JSON payload from an HTTP request into the provided structure.
//
// Parameters:
// - r: A pointer to the http.Request containing the JSON payload in the request body.
// - payload: A reference to the variable where the parsed data should be stored.
//
// Returns:
// - An error if the request body is nil or the JSON decoding fails; otherwise, nil.
func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing Body Request")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

// WriteJSON encodes a value as JSON and writes it to the HTTP response with the specified status code.
//
// Parameters:
// - w: The http.ResponseWriter used to write the JSON response.
// - status: The HTTP status code to set for the response (e.g., 200 for OK).
// - value: The data to be encoded as JSON and sent in the response body.
//
// Returns:
// - An error if encoding the value as JSON fails; otherwise, nil.
func WriteJSON(w http.ResponseWriter, status int, value any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(value)
}
