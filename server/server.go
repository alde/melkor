package server

import (
	"net/http"

	"github.com/alde/melkor"
	"github.com/alde/melkor/config"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

// NewRouter is used to create a new HTTP router
func NewRouter(cfg *config.Config, crawlers melkor.Crawlers) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	h := NewHandler(cfg, crawlers)

	for _, route := range routes(h) {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(prometheus.InstrumentHandler(route.Name, route.Handler))
	}
	return router
}

// Route enforces the structure of a route
type route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.Handler
}

func routes(h *Handler) []route {
	return []route{
		{
			Name:    "ListResources",
			Method:  "GET",
			Pattern: "/api/v1/aws/{resource}",
			Handler: h.ListAWSResources(),
		},
		{
			Name:    "GetSingleResource",
			Method:  "GET",
			Pattern: "/api/v1/aws/{resource}/{id}",
			Handler: h.GetSingleAWSResource(),
		},
		{
			Name:    "ServiceMetadata",
			Method:  "GET",
			Pattern: "/service-metadata",
			Handler: h.ServiceMetadata(),
		},
	}
}
