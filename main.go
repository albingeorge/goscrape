package main

import (
	"io/ioutil"
	"net/http"

	"github.com/albingeorge/scraper/datamanager"
	"github.com/albingeorge/scraper/logger"

	"github.com/albingeorge/scraper/config"
	"github.com/albingeorge/scraper/htmlparser"
)

func main() {
	env := config.GetInstance()
	env.LoadConfigs()

	logger.Init()

	// 1. Fetch URL from config
	// 2. Fetch html from URL
	// 3. Fetch href and other tags(like <img>) from html
	// 4. Create data index by filtering base domain name
	// 5. Start crawling the data by following href tags
	// 6. Start downloading the files from data

	url := env.GetConf().URL

	urlListObj := datamanager.GetURLListInstance()
	urlList := urlListObj.Get()

	urlListObj.Add(url)
	// fmt.Println(urlList)
	lgr := logger.Get()

	for true {
		urlList = urlListObj.Get()

		url = fetchURLToParse(urlList)

		if url != "" {
			lgr.Info("url", map[string]interface{}{"url": url})
			htmlReader, err := fetchHTMLReader(url)
			if err != nil {
				panic(err)
			}

			htmlparser.ParseHTMLString(url, htmlReader)

			// fmt.Println("URL List:")
			// for url := range urlList.Get() {
			// 	fmt.Println(url)
			// }

			// fmt.Println("\n\nResource List:")
			// for r := range resourceList.Get() {
			// 	fmt.Println(r)
			// }
		} else {
			break
		}
	}

	logURLList(urlList)
	logResourceList(datamanager.GetResourceListInstance().Get())
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

func fetchURLToParse(urlList map[string]bool) string {
	for url, done := range urlList {
		if done == false {
			return url
		}
	}

	return ""
}

func logURLList(urlList map[string]bool) {
	lgr := logger.Get()
	lgr.Info("URL List", map[string]interface{}{"urls": urlList})
}

func logResourceList(resources map[string]datamanager.Resource) {
	lgr := logger.Get()
	urls := []string{}
	for url := range resources {
		urls = append(urls, url)
	}
	lgr.Info("Resource URL", map[string]interface{}{"url": urls})
}
