package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alde/melkor"
	"github.com/alde/melkor/config"
	"github.com/alde/melkor/mock"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Test_ListAWSResources_Unknown(t *testing.T) {
	m := mux.NewRouter()
	config := &config.Config{}
	coll := melkor.Crawlers{}
	h := NewHandler(config, coll)
	m.HandleFunc("/api/v1/aws/{resource}", h.ListAWSResources())
	wr := httptest.NewRecorder()

	r, _ := http.NewRequest("GET", "/api/v1/aws/unknown", nil)
	m.ServeHTTP(wr, r)
	assert.Equal(t, wr.Code, http.StatusNotFound)
}

func setupListAWSResources(data []map[string]interface{}) (*mux.Router, *httptest.ResponseRecorder) {
	m := mux.NewRouter()
	config := &config.Config{}
	mc := &mock.Crawler{Data: data}
	coll := melkor.Crawlers{mc.Resource(): mc}
	h := NewHandler(config, coll)
	m.HandleFunc("/api/v1/aws/{resource}", h.ListAWSResources())
	wr := httptest.NewRecorder()
	return m, wr
}

func Test_ListAWSResources_Empty(t *testing.T) {
	m, wr := setupListAWSResources(mock.EmptyCrawlerData())

	r, _ := http.NewRequest("GET", "/api/v1/aws/mock", nil)
	m.ServeHTTP(wr, r)

	assert.Equal(t, wr.Code, http.StatusOK)

	expected := `[]`
	actual := strings.TrimSpace(wr.Body.String())
	assert.Equal(t, expected, actual)
}

func Test_ListAWSResources_All(t *testing.T) {
	m, wr := setupListAWSResources(mock.FullCrawlerData())

	r, _ := http.NewRequest("GET", "/api/v1/aws/mock", nil)
	m.ServeHTTP(wr, r)

	assert.Equal(t, wr.Code, http.StatusOK)
	expected := []string{"m-0", "m-1", "m-2", "m-3"}
	var actual []string
	err := json.Unmarshal(wr.Body.Bytes(), &actual)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func Test_ListAWSResources_Limit(t *testing.T) {
	m, wr := setupListAWSResources(mock.FullCrawlerData())

	r, _ := http.NewRequest("GET", "/api/v1/aws/mock?_limit=1", nil)
	m.ServeHTTP(wr, r)

	assert.Equal(t, wr.Code, http.StatusOK)

	expected := []string{"m-0"}
	var actual []string
	err := json.Unmarshal(wr.Body.Bytes(), &actual)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func Test_ListAWSResources_LimitFail(t *testing.T) {
	m, wr := setupListAWSResources(mock.FullCrawlerData())

	r, _ := http.NewRequest("GET", "/api/v1/aws/mock?_limit=one", nil)
	m.ServeHTTP(wr, r)

	assert.Equal(t, wr.Code, http.StatusBadRequest)

	expected := map[string]string{"error": "Bad limit parameter"}
	var actual map[string]string
	err := json.Unmarshal(wr.Body.Bytes(), &actual)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func Test_ListAWSResources_Expand(t *testing.T) {
	m, wr := setupListAWSResources(mock.FullCrawlerData())

	r, _ := http.NewRequest("GET", "/api/v1/aws/mock?_expand=true", nil)
	m.ServeHTTP(wr, r)

	assert.Equal(t, wr.Code, http.StatusOK)

	expected := []map[string]interface{}{
		{"id": "m-0", "name": "Mock 0", "region": "eu-west-1"},
		{"id": "m-1", "name": "Mock 1", "region": "eu-west-1"},
		{"id": "m-2", "name": "Mock 2", "region": "eu-west-1"},
		{"id": "m-3", "name": "Mock 3", "region": "eu-west-1"},
	}
	var actual []map[string]interface{}
	err := json.Unmarshal(wr.Body.Bytes(), &actual)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func setupGetSingleAWSResource() (*mux.Router, *httptest.ResponseRecorder) {
	m := mux.NewRouter()
	config := &config.Config{}
	mc := &mock.Crawler{Data: mock.FullCrawlerData()}
	coll := melkor.Crawlers{mc.Resource(): mc}
	h := NewHandler(config, coll)
	m.HandleFunc("/api/v1/aws/{resource}/{id}", h.GetSingleAWSResource())
	wr := httptest.NewRecorder()

	return m, wr
}

func Test_GetSingleAWSResource(t *testing.T) {
	m, wr := setupGetSingleAWSResource()

	r, _ := http.NewRequest("GET", "/api/v1/aws/mock/m-1", nil)
	m.ServeHTTP(wr, r)

	assert.Equal(t, wr.Code, http.StatusOK)

	expected := map[string]interface{}{
		"id": "m-1", "name": "Mock 1", "region": "eu-west-1",
	}
	var actual map[string]interface{}
	err := json.Unmarshal(wr.Body.Bytes(), &actual)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func Test_GetSingleAWSResource_NotFound(t *testing.T) {
	m, wr := setupGetSingleAWSResource()

	r, _ := http.NewRequest("GET", "/api/v1/aws/mock/m-99", nil)
	m.ServeHTTP(wr, r)

	assert.Equal(t, wr.Code, http.StatusNotFound)

	expected := map[string]interface{}{"error": "Not Found"}
	var actual map[string]interface{}
	err := json.Unmarshal(wr.Body.Bytes(), &actual)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func Test_GetSingleAWSResource_UnknownResource(t *testing.T) {
	m, wr := setupGetSingleAWSResource()

	r, _ := http.NewRequest("GET", "/api/v1/aws/unmock/m-99", nil)
	m.ServeHTTP(wr, r)

	assert.Equal(t, wr.Code, http.StatusNotFound)

	expected := map[string]interface{}{"error": "Not Found"}
	var actual map[string]interface{}
	err := json.Unmarshal(wr.Body.Bytes(), &actual)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func Test_ServiceMetadata(t *testing.T) {
	m := mux.NewRouter()
	config := &config.Config{}
	mc := &mock.Crawler{
		CountFn: func() int {
			return 20
		},
		LastCrawledFn: func() time.Time {
			return time.Now()
		},
	}
	coll := melkor.Crawlers{mc.Resource(): mc}
	h := NewHandler(config, coll)
	m.HandleFunc("/service-metadata", h.ServiceMetadata())
	wr := httptest.NewRecorder()

	r, _ := http.NewRequest("GET", "/service-metadata", nil)
	m.ServeHTTP(wr, r)

	assert.Equal(t, wr.Code, http.StatusOK)

	var actual map[string]interface{}
	err := json.Unmarshal(wr.Body.Bytes(), &actual)
	assert.Nil(t, err)

	expectedKeys := []string{
		"service_name", "service_version", "aws_region",
		"crawlers", "description", "owner",
	}

	for _, k := range expectedKeys {
		_, ok := actual[k]
		assert.True(t, ok)
	}

	crawlers := actual["crawlers"].([]interface{})
	assert.Len(t, crawlers, 1, "Mock only has one crawler")

	crawler0 := crawlers[0].(map[string]interface{})
	if v, ok := crawler0["resource"]; ok {
		assert.Equal(t, v, mc.Resource())
	}

	if v, ok := crawler0["last_crawled"]; ok {
		lc := v.(string)
		_, err := time.Parse(time.RFC3339, lc)
		assert.Nil(t, err, "Parsing the Timestamp")
	}

	if v, ok := crawler0["count"]; ok {
		assert.Equal(t, int(v.(float64)), 20)
	}
}
