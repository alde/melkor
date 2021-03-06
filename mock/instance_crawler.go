package mock

import "time"

// The InstanceCrawler mock struct holds the mocked implementation for the crawler interface
type InstanceCrawler struct {
	ResourceFn        func() string
	ResourceFnInvoked bool

	LastCrawledFn        func() time.Time
	LastCrawledFnInvoked bool

	DoCrawlFn        func() error
	DoCrawlFnInvoked bool

	ListFn        func() []string
	ListFnInvoked bool

	ListExpandedFn        func() []map[string]interface{}
	ListExpandedFnInvoked bool

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

// Resource identifies the name of the crawled resource
func (mc *InstanceCrawler) Resource() string {
	mc.ResourceFnInvoked = true
	if mc.ResourceFn == nil {
		return "Mock"
	}
	return mc.ResourceFn()
}

// LastCrawled is the timestamp of the most recent crawl
func (mc *InstanceCrawler) LastCrawled() time.Time {
	mc.LastCrawledFnInvoked = true
	return mc.LastCrawledFn()
}

// DoCrawl handles the crawling
func (mc *InstanceCrawler) DoCrawl() error {
	mc.DoCrawlFnInvoked = true
	return mc.DoCrawl()
}

// List resources
func (mc *InstanceCrawler) List() []string {
	mc.ListFnInvoked = true
	if mc.ListFn == nil {
		return mc.defaultListFn()
	}
	return mc.ListFn()
}

// ListExpanded resources
func (mc *InstanceCrawler) ListExpanded() []map[string]interface{} {
	mc.ListFnInvoked = true
	if mc.ListFn == nil {
		return mc.defaultListExpandedFn()
	}
	return mc.ListExpandedFn()
}

func (mc *InstanceCrawler) defaultListFn() []string {
	if len(mc.Data) == 0 {
		return []string{}
	}
	var data []string
	for _, d := range mc.Data {
		data = append(data, d["InstanceId"].(string))
	}
	return data
}

func (mc *InstanceCrawler) defaultListExpandedFn() []map[string]interface{} {
	var data []map[string]interface{}
	for _, d := range mc.Data {
		data = append(data, d)
	}
	return data
}

// Get returns a single resource by id
func (mc *InstanceCrawler) Get(id string) map[string]interface{} {
	mc.GetFnInvoked = true
	if mc.GetFn == nil {
		return mc.defaultGetFn(id)
	}
	return mc.GetFn(id)
}

func (mc *InstanceCrawler) defaultGetFn(id string) map[string]interface{} {
	for _, d := range mc.Data {
		if d["InstanceId"] == id {
			return d
		}
	}
	return nil
}

// Count the number of resources crawled
func (mc *InstanceCrawler) Count() int {
	mc.CountFnInvoked = true
	return mc.CountFn()
}
