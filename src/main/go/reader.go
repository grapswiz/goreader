package goreader

import (
	"code.google.com/p/go.net/html"
	"fmt"
	"goweb"
	"io/ioutil"
	"net/http"

	"appengine"
	"appengine/urlfetch"
)

type FeedsController struct{}

func (r *FeedsController) HandleRequest(c *goweb.Context) {
	context := appengine.NewContext(c.Request)
	client := urlfetch.Client(context)
	method := c.PathParams["method"]
	if method == "get" {
		feeds, err := getFeeds(context)
		if err != nil {
			fmt.Fprintf(c.ResponseWriter, "Feed Not Found")
			return
		}
		fmt.Fprintf(c.ResponseWriter, "%s", feeds)
	} else if method == "post" {
		url, err := feedUrl(client, "http://blog.alexmaccaw.com/")
		if err != nil {
			fmt.Fprintf(c.ResponseWriter, "Feed URL Not Found")
			return
		}
		err = createFeed(context, url)
		if err != nil {
			panic(err)
		}
		data := feedData(client, url)
		fmt.Fprintf(c.ResponseWriter, "%s", string(data))
	}
}

func init() {
	http.HandleFunc("/v1/auth", authHandler)
	var feedsController *FeedsController = new(FeedsController)
	goweb.Map("/v1/api/feeds/{method}", feedsController)

	goweb.ConfigureDefaultFormatters()
	http.Handle("/v1/api/", goweb.DefaultHttpHandler)
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
