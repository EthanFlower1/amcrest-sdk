package amcrest

import "fmt"

// APIError represents an error response from the Amcrest API.
type APIError struct {
	StatusCode int
	Code       int
	Message    string
}

func (e *APIError) Error() string {
	if e.Code != 0 {
		return fmt.Sprintf("amcrest: HTTP %d, error code %d: %s", e.StatusCode, e.Code, e.Message)
	}
	if e.Message != "" {
		return fmt.Sprintf("amcrest: HTTP %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("amcrest: HTTP %d", e.StatusCode)
}
