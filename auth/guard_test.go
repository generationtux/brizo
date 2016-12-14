package auth

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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
func TestAPIMiddlewareInvalidHeader(t *testing.T) {
	rw := new(responseWriterMock)
	request, _ := http.NewRequest("GET", "/api/v1/foo", nil)
	request.Header.Set("Authorization", "NotSupported token123")
	rw.On("Write", []byte("invalid authorization header")).Once()
	rw.On("WriteHeader", 401).Once()

	APIMiddleware(rw, request, http.HandlerFunc(mockHTTPHandler))

	rw.AssertExpectations(t)
}

func TestExtractJWTFromHeader(t *testing.T) {
	header := "Bearer tokenabc123"
	token := extractJWTFromHeader(header)
	assert.Equal(t, "tokenabc123", token)
}

func TestAPIMiddlewareValidJWTToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "secret")
	os.Setenv("JWT_ALGO", "HS256")
	user := User{
		Name:     "Foo Bar",
		Email:    "foo@example.com",
		Username: "foobar",
	}
	token, err := CreateJWTToken(user)
	if err != nil {
		t.Fatal(err)
	}

	rw := new(responseWriterMock)
	request, _ := http.NewRequest("GET", "/api/v1/foo", nil)
	request.Header.Set("Authorization", "Bearer "+token)
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
