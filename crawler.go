package melkor

import (
	"strings"
	"time"
)

// The Crawler interface sets up the contract for a crawler
type Crawler interface {
	DoCrawl() error
	Resource() string
	List() []string
	ListExpanded() []map[string]interface{}
	Get(id string) map[string]interface{}
	LastCrawled() time.Time
	Count() int
}

// The Crawlers struct holds all the creepy crawlies
type Crawlers map[string]Crawler

// Get fetches a Crawler case-insensitively.
func (c Crawlers) Get(r string) Crawler {
	if val, ok := c[strings.Title(strings.ToLower(r))]; ok {
		return val
	}
	return nil
}
