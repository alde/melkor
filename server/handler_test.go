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
	"github.com/alde/melkor/fixtures"
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
	mc := &mock.InstanceCrawler{Data: data}
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
	m, wr := setupListAWSResources(fixtures.FullCrawlerData(4))

	r, _ := http.NewRequest("GET", "/api/v1/aws/mock", nil)
	m.ServeHTTP(wr, r)

	assert.Equal(t, wr.Code, http.StatusOK)
	expected := []string{"i-0", "i-1", "i-2", "i-3"}
	var actual []string
	err := json.Unmarshal(wr.Body.Bytes(), &actual)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func Test_ListAWSResources_Limit(t *testing.T) {
	m, wr := setupListAWSResources(fixtures.FullCrawlerData(3))

	r, _ := http.NewRequest("GET", "/api/v1/aws/mock?_limit=1", nil)
	m.ServeHTTP(wr, r)

	assert.Equal(t, wr.Code, http.StatusOK)

	expected := []string{"i-0"}
	var actual []string
	err := json.Unmarshal(wr.Body.Bytes(), &actual)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func Test_ListAWSResources_LimitFail(t *testing.T) {
	m, wr := setupListAWSResources(fixtures.FullCrawlerData(3))

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
	m, wr := setupListAWSResources(fixtures.FullCrawlerData(3))

	r, _ := http.NewRequest("GET", "/api/v1/aws/mock?_expand=true", nil)
	m.ServeHTTP(wr, r)

	assert.Equal(t, http.StatusOK, wr.Code)

	expected := fixtures.ExpectedFullResponse(3)
	actual := strings.TrimRight(wr.Body.String(), "\n")

	assert.Equal(t, expected, actual)
}

func Test_ListAWSResources_Expand_Limit(t *testing.T) {
	m, wr := setupListAWSResources(fixtures.FullCrawlerData(3))

	r, _ := http.NewRequest("GET", "/api/v1/aws/mock?_expand=true&_limit=1", nil)
	m.ServeHTTP(wr, r)

	assert.Equal(t, http.StatusOK, wr.Code)

	var actual []map[string]interface{}
	err := json.Unmarshal(wr.Body.Bytes(), &actual)

	assert.Nil(t, err)
	assert.Len(t, actual, 1)
	for _, a := range actual {
		assert.Contains(t, a, "InstanceId")
		assert.Equal(t, a["InstanceId"], "i-0")
		assert.Contains(t, a, "Tags")
		for _, tag := range a["Tags"].([]interface{}) {
			tag0 := tag.(map[string]interface{})
			key := tag0["Key"].(string)
			value := tag0["Value"].(string)
			assert.Contains(t, tag0, key)
			assert.Equal(t, value, tag0[key])
		}
	}
}

func Test_ListAWSResources_Filter(t *testing.T) {
	m, wr := setupListAWSResources(fixtures.FullCrawlerData(3))

	r, _ := http.NewRequest("GET", "/api/v1/aws/mock?_expand=true&_filter=(Tags.Team:team2)", nil)
	m.ServeHTTP(wr, r)

	assert.Equal(t, http.StatusOK, wr.Code)

	var actual []map[string]interface{}
	err := json.Unmarshal(wr.Body.Bytes(), &actual)

	assert.Nil(t, err)
	assert.Len(t, actual, 1)
	assert.Contains(t, actual[0], "InstanceId")
	assert.Equal(t, actual[0]["InstanceId"], "i-2")
}

func Test_ListAWSResources_Filter_ValueInsensitive(t *testing.T) {
	m, wr := setupListAWSResources(fixtures.FullCrawlerData(3))

	r, _ := http.NewRequest("GET", "/api/v1/aws/mock?_expand=true&_filter=(Tags.Team:TEAM2)", nil)
	m.ServeHTTP(wr, r)

	assert.Equal(t, http.StatusOK, wr.Code)

	var actual []map[string]interface{}
	err := json.Unmarshal(wr.Body.Bytes(), &actual)

	assert.Nil(t, err)
	assert.Len(t, actual, 1)
	assert.Contains(t, actual[0], "InstanceId")
	assert.Equal(t, actual[0]["InstanceId"], "i-2")
}

func Test_ListAWSResources_Filter_InvalidFormat(t *testing.T) {
	m, wr := setupListAWSResources(fixtures.FullCrawlerData(3))

	r, _ := http.NewRequest("GET", "/api/v1/aws/mock?_expand=true&_filter=Tags.Team:TEAM2", nil)
	m.ServeHTTP(wr, r)

	assert.Equal(t, http.StatusInternalServerError, wr.Code)

	var actual map[string]interface{}
	err := json.Unmarshal(wr.Body.Bytes(), &actual)

	assert.Nil(t, err)
	assert.Contains(t, actual, "error")
	assert.Contains(t, actual["error"], "invalid format of filter")
}

func Test_ListAWSResources_Filter_And_Limit(t *testing.T) {
	m, wr := setupListAWSResources(fixtures.FullCrawlerData(3))

	r, _ := http.NewRequest("GET", "/api/v1/aws/mock?_expand=true&_filter=(Tags.Team:team2)&_limit=1", nil)
	m.ServeHTTP(wr, r)

	assert.Equal(t, http.StatusOK, wr.Code)

	var actual []map[string]interface{}
	err := json.Unmarshal(wr.Body.Bytes(), &actual)

	assert.Nil(t, err)
	assert.Len(t, actual, 1)
	assert.Contains(t, actual[0], "InstanceId")
	assert.Equal(t, actual[0]["InstanceId"], "i-2")
}

func setupGetSingleAWSResource() (*mux.Router, *httptest.ResponseRecorder) {
	m := mux.NewRouter()
	config := &config.Config{}
	mc := &mock.InstanceCrawler{Data: fixtures.FullCrawlerData(3)}
	coll := melkor.Crawlers{mc.Resource(): mc}
	h := NewHandler(config, coll)
	m.HandleFunc("/api/v1/aws/{resource}/{id}", h.GetSingleAWSResource())
	wr := httptest.NewRecorder()

	return m, wr
}

func Test_GetSingleAWSResource(t *testing.T) {
	m, wr := setupGetSingleAWSResource()

	r, _ := http.NewRequest("GET", "/api/v1/aws/mock/i-1", nil)
	m.ServeHTTP(wr, r)

	assert.Equal(t, http.StatusOK, wr.Code)

	var actual map[string]interface{}
	err := json.Unmarshal(wr.Body.Bytes(), &actual)

	assert.Nil(t, err)
	assert.Contains(t, actual, "InstanceId")
	assert.Equal(t, actual["InstanceId"], "i-1")
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
	mc := &mock.InstanceCrawler{
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
