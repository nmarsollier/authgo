package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	_ "github.com/nmarsollier/authgo/docs"
)

// Requests Test functions
func TestGetRequest(url string, tokenString string) (*http.Request, *httptest.ResponseRecorder) {
	req, _ := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
	if len(tokenString) > 0 {
		req.Header.Add("Authorization", "Bearer "+tokenString)
	}
	w := httptest.NewRecorder()
	return req, w
}

func TestPostRequest(url string, body interface{}, tokenString string) (*http.Request, *httptest.ResponseRecorder) {
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if len(tokenString) > 0 {
		req.Header.Add("Authorization", "Bearer "+tokenString)
	}
	w := httptest.NewRecorder()
	return req, w
}
