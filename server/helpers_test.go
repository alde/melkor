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

var filterTests = []struct {
	input string
	keys  []string
	value string
}{
	{"(foo.bar:baz)", []string{"foo", "bar"}, "baz"},
	{"(foo:baz)", []string{"foo"}, "baz"},
	{"(egg.bacon.ham:breakfast)", []string{"egg", "bacon", "ham"}, "breakfast"},
	{"(breakfast:bacon.ham)", []string{"breakfast"}, "bacon.ham"},
}

func Test_parseFilter(t *testing.T) {
	for _, tt := range filterTests {
		ak, av, _ := parseFilter(tt.input)
		assert.Equal(t, tt.keys, ak)
		assert.Equal(t, tt.value, av)
	}
}

var badFilters = []string{
	")(", "(foo.:)", "foo.bar:bib", "(foo:bar:baz)", "(foo.bar.baz)",
}

func Test_parseFilter_Fail(t *testing.T) {
	for _, input := range badFilters {
		_, _, err := parseFilter(input)
		assert.NotNil(t, err, "Parsing %s", input)
	}
}

func Test_deepSearch(t *testing.T) {
	input := map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": "baz",
		},
	}

	res := deepSearch(input, []string{"foo", "bar"}, "baz")
	assert.True(t, res)
}

func Test_deepSearch_Two(t *testing.T) {
	input := map[string]interface{}{
		"foo": []map[string]interface{}{
			{"bar": "baz"},
			{"bar": "bingo"},
		},
	}

	res := deepSearch(input, []string{"foo", "bar"}, "bingo")
	assert.True(t, res)
}

func Test_deepSearch_Three(t *testing.T) {
	var tags []interface{}
	tags = append(tags, map[string]interface{}{
		"Key":   "Team",
		"Value": "TestTeam",
		"Team":  "TestTeam",
	})
	input := map[string]interface{}{
		"foo": []map[string]interface{}{
			{"bar": "baz"},
			{"bar": "bingo"},
			{"tags": tags},
		},
	}
	t.Logf("%+v", input)

	res := deepSearch(input, []string{"foo", "tags", "Team"}, "TestTeam")
	assert.True(t, res)
}

func Test_applyFilter(t *testing.T) {
	filter := "(foo.bar:baz)"
	input := []map[string]interface{}{
		{
			"foo": map[string]interface{}{
				"bar": "baz",
			},
		},
		{
			"foo": map[string]interface{}{
				"bar": "bingo",
			},
		},
	}

	actual, err := applyFilter(filter, input)
	if err != nil {
		t.Error(err)
	}
	expected := []map[string]interface{}{
		{
			"foo": map[string]interface{}{"bar": "baz"},
		},
	}
	assert.Equal(t, expected, actual)
}

func Test_applyFilter_Two(t *testing.T) {
	filter := "(foo.bar:bob)"
	input := []map[string]interface{}{
		{
			"foo": []map[string]interface{}{
				{"bar": "baz"},
				{"bar": "bingo"},
				{"bar": "beef"},
			},
		},
		{
			"foo": []map[string]interface{}{
				{"bar": "baz"},
				{"bar": "bob"},
			},
		},
	}

	actual, err := applyFilter(filter, input)
	if err != nil {
		t.Error(err)
	}
	expected := []map[string]interface{}{
		{
			"foo": []map[string]interface{}{{"bar": "baz"}, {"bar": "bob"}},
		},
	}
	assert.Equal(t, expected, actual)
}
