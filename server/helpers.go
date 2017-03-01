package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

const (
	contentTypeJSON = "application/json; charset=UTF-8"
)

func writeJSON(status int, data interface{}, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func notFound(w http.ResponseWriter) error {
	return writeError(http.StatusNotFound, "Not Found", w)
}

func writeError(status int, message string, w http.ResponseWriter) error {
	data := make(map[string]string)
	data["error"] = message
	return writeJSON(status, data, w)
}

func parseLimit(r *http.Request) (int, error) {
	l := r.FormValue("_limit")
	if l == "" {
		return 0, nil
	}
	return strconv.Atoi(l)
}

func applyLimit(data []string, limit int) []string {
	if limit == 0 {
		return data
	}
	var collection []string
	for idx, el := range data {
		if idx == limit {
			break
		}
		collection = append(collection, el)
	}
	return collection
}

func applyLimitExpanded(data []map[string]interface{}, limit int) []map[string]interface{} {
	if limit == 0 {
		return data
	}
	var collection []map[string]interface{}
	for idx, el := range data {
		if idx == limit {
			break
		}
		collection = append(collection, el)
	}
	return collection
}

func applyFilter(filter string, data []map[string]interface{}) ([]map[string]interface{}, error) {
	keys, value, err := parseFilter(filter)
	if err != nil {
		return data, err
	}
	var collection []map[string]interface{}
	for i, el := range data {
		if deepSearch(el, keys, value) {
			collection = append(collection, data[i])
		}
	}

	return collection, nil
}

func deepSearch(el map[string]interface{}, keys []string, value string) bool {
	k := keys[0]
	if len(keys) == 1 {
		v := el[k]
		if v == nil {
			return false
		}

		return strings.ToLower(el[k].(string)) == strings.ToLower(value)
	}
	tail := keys[1:]

	switch el[k].(type) {
	case map[string]interface{}:
		return deepSearch(el[k].(map[string]interface{}), tail, value)
	case []map[string]interface{}:
		for _, b := range el[k].([]map[string]interface{}) {
			if deepSearch(b, tail, value) {
				return true
			}
		}
	case []interface{}:
		for _, b := range el[k].([]interface{}) {
			if deepSearch(b.(map[string]interface{}), tail, value) {
				return true
			}
		}
	}

	return false
}

func parseFilter(filter string) ([]string, string, error) {
	if !strings.HasPrefix(filter, "(") && !strings.HasSuffix(filter, ")") {
		return []string{}, "", errors.New("invalid format of filter, must be surrounded by '()'")
	}
	s := strings.Trim(filter, "()")
	if strings.Count(s, ":") != 1 {
		return []string{}, "", errors.New("invalid format of filter, only one ':' allowed")
	}
	if strings.HasPrefix(s, ".") || strings.HasSuffix(s, ".") || strings.HasPrefix(s, ":") || strings.HasSuffix(s, ":") {
		return []string{}, "", errors.New("invalid format of filter, must not start or end with '.' or ':'")
	}
	if strings.Count(s, ".:") != 0 || strings.Count(s, ":.") != 0 {
		return []string{}, "", errors.New("invalid format of filter, must not have '.' adjacent to ':'")
	}
	splits := strings.Split(s, ":")
	val := splits[1]
	k0 := splits[0]
	keys := strings.Split(k0, ".")
	return keys, val, nil
}
