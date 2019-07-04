package datamanager

import (
	"strings"
	"sync"

	"github.com/albingeorge/scraper/config"
	"github.com/albingeorge/scraper/logger"
)

/*
	This list would be used to maintain a list of URLs to be parsed
	Initially, the config URL would be added to this, so that when parsing any other URLs
	if the config URL is present, we do not need to parse it again, later on.

	Logic of it's usage goes like this:
	1. Add config URL to this list
	2. Fetch and parse html of the config URL
	3. If any a href tags are found and they confirms to the config URL, add  those URLs to this list
	4. For each URLs in this list, fetch and parse again, followed by more additions to the list
	5. Repeat step 4 untill all URLs in this list are parsed
*/
type urlList struct {
	urls map[string]bool
	mu   sync.RWMutex
}

func (u *urlList) Add(url string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if strings.HasPrefix(url, "/") {
		conf := config.GetInstance()
		url = conf.GetConf().BaseURL + url
	}
	if _, ok := u.urls[url]; !ok {
		lgr := logger.Get()
		lgr.Info("ADD_URL", map[string]interface{}{"url": url})
		u.urls[url] = false
	}
}

func (u *urlList) SetDone(url string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if _, ok := u.urls[url]; ok {
		u.urls[url] = true
	}
}

func (u *urlList) Get() map[string]bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.urls
}

var u *urlList

// GetURLListInstance ...initialises the singleton URL list
func GetURLListInstance() *urlList {
	if u == nil {
		u = &urlList{
			urls: make(map[string]bool),
		}
	}
	return u
}

// Resource ..the download struct
type Resource struct {
	URL        string
	RefererURL string
}

/*
	resourceList ... keeps a list of resources for download

	This keeps a map of string to Resource object
	The index should uniquely identify one Resource object across all the pages,
	so that we would not have to download one resourse twice while scraping
*/
type ResourceList struct {
	mu        sync.RWMutex
	Resources map[string]Resource
}

// Add ...adds a resource to the resource list
func (d *ResourceList) Add(key string, r Resource) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.Resources[key]; !ok {
		d.Resources[key] = r
	}
}

// Get ...fetches the entire resource list
func (d *ResourceList) Get() map[string]Resource {
	u.mu.Lock()
	defer u.mu.Unlock()
	return d.Resources
}

var r *ResourceList

func GetResourceListInstance() *ResourceList {
	if r == nil {
		r = &ResourceList{
			Resources: make(map[string]Resource),
		}
	}
	return r
}
