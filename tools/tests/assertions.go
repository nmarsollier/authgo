package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Assertion Functions
func AssertUnauthorized(t *testing.T, w *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	fmt.Println(result)
	assert.Equal(t, result["error"], "Unauthorized")
}

func AssertDocumentNotFound(t *testing.T, w *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusNotFound, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, "Document not found", result["error"])
}

func AssertInternalServerError(t *testing.T, w *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func AssertBadRequestError(t *testing.T, w *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
