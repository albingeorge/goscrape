package config

import "sync"

const (
	//IMAGES ... images
	IMAGES = "img"
	//VIDEOS ... videos
	VIDEOS = "VIDEOS"
)

// Environment ... singleton of this instance would be created
type Environment struct {
}

// Configs ...to represent the internal config
type Configs struct {
	URL                string
	BaseURL            string
	RelativeURL        string
	Whitelist          []string
	ConcurrentRequests int
}

var instance *Environment
var once sync.Once
var config Configs

// GetInstance ...fetches an instance of Config once
func GetInstance() *Environment {
	once.Do(func() {
		instance = &Environment{}
	})
	return instance
}

//LoadConfigs ...Loads configs from the cli args
func (e *Environment) LoadConfigs() {
	// todo: generate this based on a config file
	config = Configs{
		URL:                "https://www.mangapanda.com/one-piece/1",
		Whitelist:          []string{IMAGES},
		ConcurrentRequests: 1,
		BaseURL:            "https://www.mangapanda.com",
		RelativeURL:        "/one-piece/1",
	}
}

// GetConf ...Gets the conf object
func (e *Environment) GetConf() Configs {
	return config
}
