package server

import (
	"encoding/json"
	"net/http"
	"strconv"
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
