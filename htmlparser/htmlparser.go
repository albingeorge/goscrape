package htmlparser

import (
	"errors"
	"strings"

	"github.com/albingeorge/scraper/config"
	"github.com/albingeorge/scraper/datamanager"
	"golang.org/x/net/html"
)

// ParseHTMLString ...parses an html string input
func ParseHTMLString(url string, htmldata string) {
	r := strings.NewReader(htmldata)

	z := html.NewTokenizer(r)

	urlList := datamanager.GetURLListInstance()

	c := config.GetInstance().GetConf()
	resourceList := datamanager.GetResourceListInstance()

testLoop:
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			urlList.SetDone(url)
			break testLoop
		case tt == html.StartTagToken || tt == html.SelfClosingTagToken:
			t := z.Token()
			isAnchor := t.Data == "a"
			if isAnchor {
				// todo: manage this error here
				url, _ := getURLFromToken(t)

				if shouldAddToURLList(url, c) {
					urlList.Add(url)
				}
			}

			if shouldAddToResourceList(t, c) {
				addToResourceList(t, resourceList)
			}
		}
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

func shouldAddToURLList(url string, conf config.Configs) bool {
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
	for _, b := range conf.Whitelist {
		if b == t.Data {
			return true
		}
	}
	return false
}

func addToResourceList(t html.Token, r *datamanager.ResourceList) {
	switch t.Data {
	case "img":
		key, resource := generateImageResource(t)
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
