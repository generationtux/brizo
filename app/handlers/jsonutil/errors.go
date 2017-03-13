package jsonutil

import (
	"encoding/json"
	"net/http"
)

// SelfRespondingJSONError is an interfaces representing the ability to render
// its own contents.
type SelfRespondingJSONError interface {
	// Render will write the contents of the JSON error to the provided
	// http.ResponseWriter
	Render(w http.ResponseWriter)
}

// GenericSelfRespondingJSONError is a generic implementation of of the
// SelfRespondingJSONError interface. This is good to use when there is little
// to no logic in defining error messages or status codes.
type GenericSelfRespondingJSONError struct {
	// Details is a human readable message that describes the error.
	Details string `json:"details,omitempty"`
	// Status is appropriate the HTTP status code for the error.
	Status int `json:"status,string"`
}

// Render writes the contents of the GenericSelfRespondingJSONError to the
// provided http.ResponseWriter.
func (srje *GenericSelfRespondingJSONError) Render(w http.ResponseWriter) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(srje.Status)
	payload, _ := json.Marshal(srje)
	w.Write(payload)
}

// DatabaseConnectError is a wrapper for a new GenericSelfRespondingJSONError
// that will provide a predefined message for a database connection error. It
// returns a SelfRespondingJSONError, in the form of a
// GenericSelfRespondingJSONError, so that rendering can be chained.
//
// Example:
//   DatabaseConnectError().Render(w)
func DatabaseConnectError() SelfRespondingJSONError {
	return &GenericSelfRespondingJSONError{
		Details: "there was an error when attempting to connect to the database",
		Status:  http.StatusInternalServerError,
	}
}

// KubeClientConnectionError provides a predefined message for a Kubernetes
// client connection error. Rendering can be chained.
func KubeClientConnectionError() SelfRespondingJSONError {
	return &GenericSelfRespondingJSONError{
		Details: "unable to reach Kubernetes",
		Status:  http.StatusServiceUnavailable,
	}
}
