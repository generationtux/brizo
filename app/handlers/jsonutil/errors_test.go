package jsonutil

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SpyResponseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (w *SpyResponseWriter) Status() int {
	return w.status
}

func (w *SpyResponseWriter) Write(p []byte) (n int, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(p)
}

func (w *SpyResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	// Check after in case there's error handling in the wrapped ResponseWriter.
	if w.wroteHeader {
		return
	}
	w.status = code
	w.wroteHeader = true
}

func assertImplementsInterface(t *testing.T, srje SelfRespondingJSONError) {
	assert.Implements(t, (*SelfRespondingJSONError)(nil), srje)
}

func assertErrorContainsContents(t *testing.T, srje SelfRespondingJSONError, expPayload string, expStatus uint) {
	w := httptest.NewRecorder()
	srje.Render(w)
	payload, _ := ioutil.ReadAll(w.Result().Body)
	assert.Contains(t, string(payload), expPayload)
	assert.Contains(t, string(payload), strconv.FormatUint(uint64(expStatus), 10))
}

func TestErrorTypes(t *testing.T) {
	dce := DatabaseConnectError()
	assertImplementsInterface(t, dce)
	assertErrorContainsContents(t,
		dce,
		"there was an error when attempting to connect to the database",
		http.StatusInternalServerError)

	kcce := KubeClientConnectionError()
	assertImplementsInterface(t, kcce)
	assertErrorContainsContents(t,
		kcce,
		"unable to reach Kubernetes",
		http.StatusServiceUnavailable)
}
