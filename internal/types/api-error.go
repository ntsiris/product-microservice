package types

import "fmt"

type APIError struct {
	Code          int    `json:"code"`
	Message       string `json:"message"`
	Operation     string `json:"operation"`
	EmbeddedError string `json:"embeddedError"`
}

func (err *APIError) Error() string {
	return ""
}

func FormatOperation(method, path string) string {
	return fmt.Sprintf("%s %s", method, path)
}
