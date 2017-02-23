package mock

import (
	"fmt"
	"time"
)

// The Crawler mock struct holds the mocked implementation for the crawler interface
type Crawler struct {
	ResourceFn        func() string
	ResourceFnInvoked bool

	LastCrawledFn        func() time.Time
	LastCrawledFnInvoked bool

	DoCrawlFn        func() error
	DoCrawlFnInvoked bool

	ListFn        func(int, bool) interface{}
	ListFnInvoked bool

	GetFn        func(string) map[string]interface{}
	GetFnInvoked bool

	CountFn        func() int
	CountFnInvoked bool

	Data []map[string]interface{}
}

// EmptyCrawlerData is an empty mock
func EmptyCrawlerData() []map[string]interface{} {
	return []map[string]interface{}{}
}

// FullCrawlerData is a mock with some data
func FullCrawlerData() []map[string]interface{} {
	var data []map[string]interface{}
	for i := 0; i <= 3; i++ {
		inner := make(map[string]interface{})

		inner["id"] = fmt.Sprintf("m-%d", i)
		inner["name"] = fmt.Sprintf("Mock %d", i)
		inner["region"] = "eu-west-1"

		data = append(data, inner)
	}
	return data
}

// Resource identifies the name of the crawled resource
func (mc *Crawler) Resource() string {
	mc.ResourceFnInvoked = true
	if mc.ResourceFn == nil {
		return "Mock"
	}
	return mc.ResourceFn()
}

// LastCrawled is the timestamp of the most recent crawl
func (mc *Crawler) LastCrawled() time.Time {
	mc.LastCrawledFnInvoked = true
	return mc.LastCrawledFn()
}

// DoCrawl handles the crawling
func (mc *Crawler) DoCrawl() error {
	mc.DoCrawlFnInvoked = true
	return mc.DoCrawl()
}

// List resources
func (mc *Crawler) List(limit int, expand bool) interface{} {
	mc.ListFnInvoked = true
	if mc.ListFn == nil {
		return mc.DefaultListFn(limit, expand)
	}
	return mc.ListFn(limit, expand)
}

func (mc *Crawler) DefaultListFn(limit int, expand bool) interface{} {
	if len(mc.Data) == 0 {
		return []string{}
	}
	if expand {
		var data []map[string]interface{}
		for idx, d := range mc.Data {
			if limit > 0 && idx == limit {
				break
			}
			data = append(data, d)
		}
		return data
	}

	var data []string
	for idx, d := range mc.Data {
		if limit > 0 && idx == limit {
			break
		}
		data = append(data, d["id"].(string))
	}
	return data
}

// Get returns a single resource by id
func (mc *Crawler) Get(id string) map[string]interface{} {
	mc.GetFnInvoked = true
	if mc.GetFn == nil {
		return mc.DefaultGetFn(id)
	}
	return mc.GetFn(id)
}

func (mc *Crawler) DefaultGetFn(id string) map[string]interface{} {
	for _, d := range mc.Data {
		if d["id"] == id {
			return d
		}
	}
	return nil
}

// Count the number of resources crawled
func (mc *Crawler) Count() int {
	mc.CountFnInvoked = true
	return mc.CountFn()
}
