package server

import (
	"testing"

	"github.com/alde/melkor"
	"github.com/alde/melkor/config"

	"github.com/stretchr/testify/assert"
)

var (
	cfg = &config.Config{}
	crw = melkor.Crawlers{}
)

func Test_NewRouter(t *testing.T) {
	h := NewHandler(cfg, crw)
	nr := NewRouter(cfg, crw)

	for _, r := range routes(h) {
		assert.NotNil(t, nr.GetRoute(r.Name))
	}
}

func Test_routes(t *testing.T) {
	h := NewHandler(cfg, crw)
	assert.Len(t, routes(h), 3, "3 routes is the magic number.")
}
