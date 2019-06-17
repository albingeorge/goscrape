package htmlparser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/albingeorge/scraper/config"
	"github.com/albingeorge/scraper/datamanager"
	"golang.org/x/net/html"
)

// ParseHTMLString ...parses an html string input
func ParseHTMLString(htmldata string) {
	r := strings.NewReader(htmldata)

	z := html.NewTokenizer(r)

	i := 1
	urlList := datamanager.GetURLListInstance()

	c := config.GetInstance().GetConf()
	resourceList := datamanager.GetResourceListInstance()

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

				if isRequiredToAdd(url, c) {
					urlList.Add(url)
				}
			}

			if shouldAddToResourceList(t, c) {
				fmt.Println("shouldAddToResourceList: ", shouldAddToResourceList(t, c))
				addToResourceList(t, resourceList)
			}
		}
	}
	// fmt.Println("URL List:")
	// for url := range urlList.Get() {
	// 	fmt.Println(url)
	// }

	fmt.Println("\n\nResource List:")
	for r := range resourceList.Get() {
		fmt.Println(r)
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

func isRequiredToAdd(url string, conf config.Configs) bool {
	// if it's a relative URL
	if strings.HasPrefix(url, "/") {
		/*
			1. Generate relative URL from config
				If global URL is "http://abc.com/a", we should get "/a"
			2. Check if the current URL starts with the relative global URL

		*/
		if strings.HasPrefix(url, conf.RelativeURL) {
			return true
		}
	} else {
		/*
			Check if the current URL is an extension of the input URL

			i.e, if the current URL is "aconfc.com/a/123" and the base URL is
			"abc.com/a", then it needs to be added.

			if the current URL is "abc.com/a/123" and the base URL is
			"abc.com/b", then it should not to be added
		*/
		if strings.HasPrefix(url, conf.URL) {
			return true
		}
	}
	return false
}

func shouldAddToResourceList(t html.Token, conf config.Configs) bool {
	fmt.Println(t.Data)
	for _, b := range conf.Whitelist {
		if b == t.Data {
			return true
		}
	}
	return false
}

func addToResourceList(t html.Token, r *datamanager.ResourceList) {
	fmt.Println("addToResourceList: ", t.Data)
	switch t.Data {
	case "img":
		key, resource := generateImageResource(t)
		fmt.Println("Adding to resource list: ", key)
		r.Add(key, resource)
	}
}

func generateImageResource(t html.Token) (string, datamanager.Resource) {
	r := datamanager.Resource{}
	for _, attr := range t.Attr {
		if attr.Key == "src" {
			r.URL = attr.Val
		}
	}
	return r.URL, r
}
