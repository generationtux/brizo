package jsonutil

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonResponseErrorInitialization(t *testing.T) {
	jre := NewJSONResponseError(500)
	assert.Equal(t, "*jsonutil.JSONResponseError", reflect.TypeOf(jre).String())
	assert.Equal(t, 500, jre.Status())
}

func TestAddingErrors(t *testing.T) {
	// Handling addition of error w/ default status code
	jre := NewJSONResponseError(500).AddError("Example error")
	error := jre.Errors()[0]
	assert.Equal(t, 500, jre.Status())
	assert.Equal(t, "Example error", error.Details)
	assert.Equal(t, 500, error.Status)

	jre = NewJSONResponseError(400).AddError("newman.gif", 401)
	error = jre.Errors()[0]
	assert.Equal(t, 400, jre.Status())
	assert.Equal(t, "newman.gif", error.Details)
	assert.Equal(t, 401, error.Status)

	jre = NewJSONResponseError(400).AddError("newman.gif", 401).AddError("User not found", 404)
	firstError := jre.Errors()[0]
	secondError := jre.Errors()[1]
	assert.Equal(t, 400, jre.Status())
	assert.Equal(t, "newman.gif", firstError.Details)
	assert.Equal(t, 401, firstError.Status)
	assert.Equal(t, "User not found", secondError.Details)
	assert.Equal(t, 404, secondError.Status)

	jre = NewJSONResponseError(400).AddError("newman.gif").AddError("User not found", 404)
	firstError = jre.Errors()[0]
	secondError = jre.Errors()[1]
	assert.Equal(t, 400, jre.Status())
	assert.Equal(t, "newman.gif", firstError.Details)
	assert.Equal(t, 400, firstError.Status)
	assert.Equal(t, "User not found", secondError.Details)
	assert.Equal(t, 404, secondError.Status)
}

func TestAddingErrorsShorthand(t *testing.T) {
	// Handling addition of error w/ default status code using shorthand style
	jre := NewJSONResponseError(500, "Example error")
	error := jre.Errors()[0]
	assert.Equal(t, 500, jre.Status())
	assert.Equal(t, "Example error", error.Details)
	assert.Equal(t, 500, error.Status)

	jre = NewJSONResponseError(400, "newman.gif", "client-side error occurred")
	firstError := jre.Errors()[0]
	secondError := jre.Errors()[1]
	assert.Equal(t, 400, jre.Status())
	assert.Equal(t, "newman.gif", firstError.Details)
	assert.Equal(t, 400, firstError.Status)
	assert.Equal(t, "client-side error occurred", secondError.Details)
	assert.Equal(t, 400, secondError.Status)
}

func TestRenderOutput(t *testing.T) {
	json := NewJSONResponseError(500).AddError("Example error").Render()
	expectedJSON := "{\"errors\":[{\"details\":\"Example error\",\"status\":\"500\"}],\"status\":\"500\"}"
	assert.Equal(t, expectedJSON, json)

	json = NewJSONResponseError(500).AddError("A different error", 503).Render()
	expectedJSON = "{\"errors\":[{\"details\":\"A different error\",\"status\":\"503\"}],\"status\":\"500\"}"
	assert.Equal(t, expectedJSON, json)
}
