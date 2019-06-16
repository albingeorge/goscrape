package htmlparser

import (
	"github.com/albingeorge/scraper/datamanager"
	// "bytes"
	"errors"
	"fmt"

	"github.com/albingeorge/scraper/config"
	"golang.org/x/net/html"

	// "io"
	"strings"
)

func ParseHtmlString(htmldata string) {
	r := strings.NewReader(htmldata)

	z := html.NewTokenizer(r)

	i := 1
	urlList := datamanager.URLList()

testLoop:
	for {
		tt := z.Next()
		i = i + 1
		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			fmt.Println("Breaking")
			break testLoop
		case tt == html.StartTagToken:
			t := z.Token()
			isAnchor := t.Data == "a"
			if isAnchor {
				// todo: manage this error here
				url, _ := getURLFromToken(t)

				if isRequiredToAdd(url) {
					fmt.Println("Adding url: ", url)
					urlList.AddURL(url)
				}
			}
		}
	}
	fmt.Println("URL List:")
	for url := range urlList.GetURLs() {
		fmt.Println(url)
	}

}

func getURLFromToken(t html.Token) (string, error) {
	for _, attr := range t.Attr {
		if attr.Key == "href" {
			return attr.Val, nil
		}
	}
	return "", errors.New("Attribute href not founnd in a tag")
}

func isRequiredToAdd(url string) bool {
	c := config.GetInstance().GetConf()

	// if it's a relative URL
	if strings.HasPrefix(url, "/") {
		/*
			1. Generate relative URL from config
				If global URL is "http://abc.com/a", we should get "/a"
			2. Check if the current URL starts with the relative global URL

		*/
		if strings.HasPrefix(url, c.RelativeURL) {
			return true
		}
	} else {
		/*
			Check if the current URL is an extension of the input URL

			i.e, if the current URL is "abc.com/a/123" and the base URL is
			"abc.com/a", then it needs to be added.

			if the current URL is "abc.com/a/123" and the base URL is
			"abc.com/b", then it should not to be added
		*/
		if strings.HasPrefix(url, c.URL) {
			return true
		}
	}
	return false
}
