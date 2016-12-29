// Package jsonutil provides support functions for dealing with json responses
// and specificially focuses on handling errors.
package jsonutil

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ResponseError holds details and a status code related to a specific error.
type ResponseError struct {
	Details string `json:"details"`
	Status  int    `json:"status,string"`
}

// ResponseErrorBag handles the management and usage of ResponseErrors. It is
// also responsible for maintaining a high-level status code, which expected to
// be reflected in the response header.
type ResponseErrorBag interface {
	// AddError will create and add a new ResponseError into this ResponseErrorBag
	// instance. It also accepts an optional status code specific to the error, or
	// defaults to the ResponseErrorBag's high-level status code.
	AddError(details string, status ...int) *ResponseErrorBag
	Errors() []ResponseError
	Render() string
	Status() int
}

// JSONResponseError implements a JSON specific ResponseErrorBag.
type JSONResponseError struct {
	bag    []ResponseError
	Output string
	status int
}

// NewJSONResponseError creates a new instance of a JSONResponseError. While a
// status is required, error messages are optionally added in a shorthand format
// if desired.
//     NewJSONResponseError(500)                   // new empty JSONResponseError
//     NewJSONResponseError(500, "internal error") // new JSONResponseError with single error
func NewJSONResponseError(status int, messages ...string) *JSONResponseError {
	jre := JSONResponseError{
		bag:    make([]ResponseError, 0), // marshals to "[]" instead of null when empty
		status: status,
	}
	if len(messages) > 0 {
		for _, message := range messages {
			jre.AddError(message, status)
		}
	}
	return &jre
}

// Errors acts as a getter to provide a slice of all of the ResponseErrors.
func (jre *JSONResponseError) Errors() []ResponseError {
	return jre.bag
}

// Status acts as a getter to provide the high-level status.
func (jre *JSONResponseError) Status() int {
	return jre.status
}

// AddError adds a new ResponseError, while accepting an optional status code
// specific to the new ResponseError.
func (jre *JSONResponseError) AddError(details string, status ...int) *JSONResponseError {
	var derivedStatus int
	if len(status) >= 1 {
		derivedStatus = status[0]
	} else {
		derivedStatus = jre.Status()
	}
	jre.bag = append(jre.bag, ResponseError{
		Details: details,
		Status:  derivedStatus,
	})

	return jre
}

// Render marshals the JSONResponseError to a valid JSON string and sets the
// Output field to the rendered output. This method does not simply marshal the
// instance, but also formats the output to a usable format for responses.
func (jre *JSONResponseError) Render() string {
	jsonBag := struct {
		Errors []ResponseError `json:"errors"`
		Status int             `json:"status,string"`
	}{
		jre.Errors(),
		jre.Status(),
	}
	response, _ := json.Marshal(jsonBag)
	jre.Output = string(response)

	return jre.Output
}

// RespondJSONError will bootstrap and write out a valid json repsonse using the
// provided JSONResponseError. It returns the results of fmt.Printf against the
// provided http.ResponseWriter, which is expected to be the number of bytes
// written and any write error encountered.
func RespondJSONError(w http.ResponseWriter, jre *JSONResponseError) (int, error) {
	w.Header().Set("content-type", "application/json")
	return fmt.Fprint(w, jre.Render())
}
