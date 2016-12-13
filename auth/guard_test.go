package auth

import (
	"net/http"
	"testing"

	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 401 if no Authorization header is present
func TestAPIMiddlewareNoHeader(t *testing.T) {
	rw := new(responseWriterMock)
	request, _ := http.NewRequest("GET", "/api/v1/foo", nil)
	rw.On("Write", []byte("not authorized")).Once()
	rw.On("WriteHeader", 401).Once()

	APIMiddleware(rw, request, http.HandlerFunc(mockHTTPHandler))

	rw.AssertExpectations(t)
}

// header must be Authorization: Bearer xxx
func TestAPIMiddlewareInalidHeader(t *testing.T) {
	rw := new(responseWriterMock)
	request, _ := http.NewRequest("GET", "/api/v1/foo", nil)
	request.Header.Set("Authorization", "NotSupported token123")
	rw.On("Write", []byte("invalid authorization header")).Once()
	rw.On("WriteHeader", 401).Once()

	APIMiddleware(rw, request, http.HandlerFunc(mockHTTPHandler))

	rw.AssertExpectations(t)
}

func TestAPIMiddlewareValidJWTToken(t *testing.T) {
	rw := new(responseWriterMock)
	request, _ := http.NewRequest("GET", "/api/v1/foo", nil)
	request.Header.Set("Authorization", "Bearer token123")
	rw.On("Write", []byte("ok")).Once()
	rw.On("WriteHeader", 200).Once()

	APIMiddleware(rw, request, http.HandlerFunc(mockHTTPHandler))

	rw.AssertExpectations(t)
}

// mocks

func mockHTTPHandler(rw http.ResponseWriter, request *http.Request) {
	rw.Write([]byte("ok"))
	rw.WriteHeader(200)
}

type responseWriterMock struct {
	mock.Mock
}

func (r *responseWriterMock) Header() http.Header {
	return http.Header{}
}
func (r *responseWriterMock) Write(content []byte) (int, error) {
	r.Called(content)
	return 0, nil
}
func (r *responseWriterMock) WriteHeader(status int) {
	r.Called(status)
}
