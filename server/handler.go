package server

import (
	"net/http"

	"github.com/alde/melkor"
	"github.com/alde/melkor/config"
	"github.com/alde/melkor/version"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Handler holds the server context
type Handler struct {
	config   *config.Config
	crawlers melkor.Crawlers
}

// NewHandler createss a new HTTP handler
func NewHandler(cfg *config.Config, crawlers melkor.Crawlers) *Handler {
	return &Handler{config: cfg, crawlers: crawlers}
}

// ListAWSResources returns a list of the requested resources, or a 404 if none
// can be found in storage
func (h *Handler) ListAWSResources() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		resource := vars["resource"]
		limit, err := parseLimit(r)
		if err != nil {
			writeError(http.StatusBadRequest, "Bad limit parameter", w)
			return
		}

		crawler := h.crawlers.Get(resource)
		if crawler == nil {
			logrus.WithField("resource", resource).Debug("Not Found")
			notFound(w)
			return
		}

		expand := r.FormValue("_expand") == "true"
		logrus.WithFields(logrus.Fields{"resource": resource, "limit": limit, "expand": expand}).Debug("Listing resources")
		if expand {
			data := crawler.ListExpanded()

			data = applyLimitExpanded(data, limit)
			writeJSON(http.StatusOK, data, w)
			return
		}
		data := crawler.List()
		data = applyLimit(data, limit)
		writeJSON(http.StatusOK, data, w)
	}
}

// GetSingleAWSResource handles fetching a single item
func (h *Handler) GetSingleAWSResource() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		resource := vars["resource"]
		id := vars["id"]
		crawler := h.crawlers.Get(resource)
		if crawler == nil {
			logrus.WithField("resource", resource).Debug("Not Found")
			notFound(w)
			return
		}
		logrus.WithFields(logrus.Fields{"resource": resource, "id": id}).Debug("Fetching single resource")
		data := crawler.Get(id)
		if data == nil {
			logrus.WithFields(logrus.Fields{"resource": resource, "id": id}).Debug("Not Found")
			notFound(w)
			return
		}

		writeJSON(http.StatusOK, data, w)
	}
}

// ServiceMetadata displays hopefully useful information about the service
func (h *Handler) ServiceMetadata() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]interface{})
		var crawlers []map[string]interface{}
		for _, c := range h.crawlers {
			inner := make(map[string]interface{})
			inner["resource"] = c.Resource()
			inner["last_crawled"] = c.LastCrawled()
			inner["count"] = c.Count()
			crawlers = append(crawlers, inner)
		}
		data["owner"] = h.config.Owner
		data["description"] = "AWS caching layer"
		data["service_name"] = "melkor"
		data["service_version"] = version.Version
		data["aws_region"] = h.config.AWSRegion
		data["crawlers"] = crawlers

		writeJSON(http.StatusOK, data, w)
	}
}
