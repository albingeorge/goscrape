package datamanager

import "sync"

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

func (u *urlList) AddURL(url string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if _, ok := u.urls[url]; !ok {
		u.urls[url] = false
	}
}

func (u *urlList) GetURLs() map[string]bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.urls
}

var u *urlList

// URLList ...initialises the URL list
func URLList() *urlList {
	if u == nil {
		u = &urlList{
			urls: make(map[string]bool),
		}
	}
	return u
}
