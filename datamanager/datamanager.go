package datamanager

import "sync"

type urlList struct {
	urls map[string]bool
	mu   sync.RWMutex
}

func (u *urlList) Add(url string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if _, ok := u.urls[url]; !ok {
		u.urls[url] = false
	}
}
