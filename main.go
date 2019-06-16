package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/albingeorge/scraper/config"
	"github.com/albingeorge/scraper/htmlparser"
)

func main() {
	env := config.GetInstance()
	env.LoadConfigs()

	// 1. Fetch URL from config
	// 2. Fetch html from URL
	// 3. Fetch href and other tags(like <img>) from html
	// 4. Create data index by filtering base domain name
	// 5. Start crawling the data by following href tags
	// 6. Start downloading the files from data

	url := env.GetConf().URL
	fmt.Println(url)

	htmlReader, err := fetchHTMLReader(url)
	if err != nil {
		panic(err)
	}

	htmlparser.ParseHTMLString(htmlReader)
}

func fetchHTMLReader(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	bytes, _ := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()
	return string(bytes), nil
}
