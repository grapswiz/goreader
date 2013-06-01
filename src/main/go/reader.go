package goreader

import (
	"code.google.com/p/go.net/html"
	"fmt"
	"io/ioutil"
	"net/http"

	"appengine"
	"appengine/urlfetch"
)

func feedHandler(rw http.ResponseWriter, req *http.Request) {
	c := appengine.NewContext(req)
	client := urlfetch.Client(c)
	method := req.Method
	if method == "GET" {
		fmt.Fprintf(rw, "get")
		url, err := feedUrl(client, "http://blog.alexmaccaw.com/")
		if err != nil {
			fmt.Fprintf(rw, "Feed URL Not Found")
			return
		}
		data := feedData(client, url)
		fmt.Fprintf(rw, "%s", string(data))

	} else if method == "POST" {
		_ = createFeed(c, req.FormValue("url"))
	}
}

func init() {
	http.HandleFunc("/v1/auth", authHandler)
	http.HandleFunc("/v1/api/feeds", feedHandler)
}

func feedUrl(client *http.Client, url string) (feedUrl string, err error) {
	feedUrl = ""

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}

	ct := resp.Header.Get("Content-Type")
	if isFeedContentType(ct) {
		return url, nil
	}

	root, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}
	d := NewDocument(root, resp.Request.URL)
	d.Find("link").Each(func(i int, s *Selection) {
		ct, exists := s.Attr("type")
		if !exists {
			return
		}
		if isFeedContentType(ct) {
			feedUrl, exists = s.Attr("href")
			return
		}
	})
	return feedUrl, nil
}

func feedData(client *http.Client, url string) []byte {
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	b, _ := ioutil.ReadAll(resp.Body)

	return b
}

func isFeedContentType(ct string) bool {
	feedContentType := []string{"application/x.atom+xml", "application/atom+xml", "application/xml", "text/xml", "application/rss+xml", "application/rdf+xml"}
	for _, contentType := range feedContentType {
		if contentType == ct {
			return true
		}
	}
	return false
}
