package api

import (
	"encoding/json"
	"errors"
	"net/http"
)

// jsonResponse writes a JSON HTTP response
func jsonResponse(rw http.ResponseWriter, data interface{}, status int) {
	content, _ := json.Marshal(data)
	rw.WriteHeader(status)
	rw.Header().Set("content-type", "application/json")
	rw.Write(content)
}

// jsonErrorResponse writes a JSON error response
func jsonErrorResponse(rw http.ResponseWriter, err string, status int) {
	data := map[string]string{
		"error": err,
	}
	jsonResponse(rw, data, status)
}

// jsonRequest unmarshals the request body to provided interface
func jsonRequest(r *http.Request, container interface{}) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return errors.New("Only JSON requests supported. Please set the Content-Type header.")
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(container)
	defer r.Body.Close()
	if err != nil {
		return err
	}

	return nil
}
