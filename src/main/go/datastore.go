package goreader

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"time"
)

type Feed struct {
	Url       string
	CreatedAt time.Time
}

type Feeds struct {
	feeds []Feed
}

type Json interface {
	toJson()
}

func (feed Feed) toJson() (string, error) {
	b, err := json.Marshal(feed)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (feeds Feeds) toJson() (string, error) {
	b, err := json.Marshal(feeds)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func getFeeds(c appengine.Context) (string, error) {
	q := datastore.NewQuery("Feed")
	feedArray := make([]Feed, 0, 10)
	_, err := q.GetAll(c, &feedArray)
	if err != nil {
		return "", nil
	}
	feeds := &Feeds{feedArray}
	json, _ := feeds.toJson()
	return json, nil
}

func createFeed(c appengine.Context, url string) error {
	feed := Feed{url, time.Now()}
	_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Feed", nil), &feed)

	return err
}
