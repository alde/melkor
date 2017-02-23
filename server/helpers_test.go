package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_writeJSON(t *testing.T) {
	wr := httptest.NewRecorder()

	writeJSON(200, "foo", wr)
	assert.Equal(t, contentTypeJSON, wr.HeaderMap["Content-Type"][0])
	assert.Equal(t, http.StatusOK, wr.Code)
}

func Test_notFound(t *testing.T) {
	wr := httptest.NewRecorder()

	notFound(wr)
	assert.Equal(t, contentTypeJSON, wr.HeaderMap["Content-Type"][0])
	assert.Equal(t, http.StatusNotFound, wr.Code)
}

func Test_writeError(t *testing.T) {
	wr := httptest.NewRecorder()

	writeError(http.StatusInternalServerError, "An error", wr)

	assert.Equal(t, contentTypeJSON, wr.HeaderMap["Content-Type"][0])
	assert.Equal(t, http.StatusInternalServerError, wr.Code)
}

func Test_parseLimit_Valid(t *testing.T) {
	expected := 1
	r, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/aws/mock/m-1?_limit=%d", expected), nil)

	actual, err := parseLimit(r)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func Test_parseLimit_Missing(t *testing.T) {
	r, _ := http.NewRequest("GET", "/api/v1/aws/mock/m-1", nil)

	actual, err := parseLimit(r)
	assert.Nil(t, err)
	assert.Equal(t, 0, actual)
}

func Test_parseLimit_Invalid(t *testing.T) {
	r, _ := http.NewRequest("GET", "/api/v1/aws/mock/m-1?_limit=one", nil)

	_, err := parseLimit(r)
	assert.NotNil(t, err)
}
