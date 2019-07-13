package htmlparser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/albingeorge/scraper/config"
	"github.com/albingeorge/scraper/datamanager"
	"github.com/albingeorge/scraper/logger"
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
				addToResourceList(t, resourceList, url)
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

func addToResourceList(t html.Token, r *datamanager.ResourceList, url string) {
	switch t.Data {
	case "img":
		resource := generateImageResource(t, url)
		r.Add(url, resource)
	}
}

func generateImageResource(t html.Token, url string) datamanager.Resource {
	directory, fileName := generateDirectoryAndFileName(url)

	r := datamanager.Resource{
		RefererURL: url,
		Directory:  directory,
		FileName:   fileName,
	}

	for _, attr := range t.Attr {
		if attr.Key == "src" {
			r.URL = attr.Val
		}
	}
	return r
}

func generateDirectoryAndFileName(url string) (string, string) {
	lgr := logger.Get()
	c := config.GetInstance().GetConf()

	/*
		Sample URL formats can be of 2 types:
		1. https://www.mangapanda.com/one-piece/940
		2. https://www.mangapanda.com/one-piece/940/5

		1. Split based on manga name("one-piece")
		2. Split the second part again based on "/" and check the size of split strings
		3. If size is 2, name the file as "001"
		4. Else if size is 3, name the file as splitString[2]
	*/

	splitByName := strings.Split(url, c.MangaName)
	splitString := strings.Split(splitByName[1], "/")

	/*
		splitString formats:
		"/940"
		"/940/2"

		In either of the cases, upon splitting by "/", the minimum length is 2
	*/
	if len(splitString) == 2 {
		return splitString[1], "001"
	} else if len(splitString) == 3 {
		pageNumber, err := strconv.Atoi(splitString[2])

		if err != nil {
			data := map[string]interface{}{
				"input_string": splitString[1],
			}
			lgr.Error("PAGE_NUMBER_FETCH_ERROR", data, "error converting string to int")
			return "", ""
		}
		pageString := fmt.Sprintf("%03d", pageNumber)

		return splitString[1], pageString
	}

	data := map[string]interface{}{
		"url": url,
	}
	lgr.Error("FILE_NAME_FETCH_ERROR", data, "error fetching file name and directory name")
	return "", ""
}
