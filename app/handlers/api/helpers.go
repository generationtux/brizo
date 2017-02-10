package api

import (
	"encoding/json"
	"net/http"
)

// jsonResponse writes a JSON HTTP response
func jsonResponse(rw http.ResponseWriter, data interface{}, status int) {
	content, _ := json.Marshal(data)
	rw.WriteHeader(status)
	rw.Header().Set("content-type", "application/json")
	rw.Write(content)
}
