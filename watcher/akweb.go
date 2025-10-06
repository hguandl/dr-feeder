package watcher

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/hguandl/dr-feeder/v2/common"
	bolt "go.etcd.io/bbolt"
)

type akWebWatcher struct {
	name     string
	category string
	debugURL string
	db       *bolt.DB
}

// NewAkWebWatcher creates a Watcher of Arknights game announcements.
func NewAkWebWatcher(dbPath string, debugURL string) (Watcher, error) {
	var err error = nil

	watcher := new(akWebWatcher)
	watcher.name = "明日方舟官网新闻"
	watcher.category = "NEWS"
	watcher.debugURL = debugURL

	watcher.db, err = bolt.Open(dbPath, 0666, nil)
	if err != nil {
		return watcher, err
	}

	_, err = watcher.update()
	return watcher, err
}

func (watcher akWebWatcher) apiURL(category string, page int) string {
	if watcher.debugURL != "" {
		return watcher.debugURL
	}
	return fmt.Sprintf("https://ak.hypergryph.com/api/news?category=%s&page=%d", category, page)
}

func (watcher akWebWatcher) fetchAPI(page int) (akWebNewsList, error) {
	var err error = nil
	var data akWebAPIPayload

	c := colly.NewCollector(
		colly.UserAgent(safariUA),
	)

	c.OnError(func(_ *colly.Response, e error) {
		err = e
	})

	c.OnResponse(func(r *colly.Response) {
		err = json.Unmarshal(r.Body, &data)
	})

	c.Visit(watcher.apiURL(watcher.category, page))
	c.Wait()
	return data.Data, err
}

func (watcher akWebWatcher) update() ([]akWebNewsData, error) {
	updates := []akWebNewsData{}

	data, err := watcher.fetchAPI(1)
	if err != nil {
		return updates, err
	}

	err = watcher.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("AkWeb"))
		for _, news := range data.List {
			v := b.Get([]byte(news.Cid))
			if v == nil {
				updates = append(updates, news)
				err = b.Put([]byte(news.Cid), []byte(news.Title))
			}
		}
		return err
	})

	return updates, err
}

func (watcher akWebWatcher) Produce(ch chan common.NotifyPayload) {
	updates, err := watcher.update()
	if err != nil {
		log.Printf("Error updating \"%s\": %v", watcher.name, err)
		return
	}

	if len(updates) == 0 {
		log.Printf("Waiting for post \"%s\"...\n", watcher.name)
	}

	for _, news := range updates {
		if !strings.Contains(news.Title, "制作组通讯") {
			log.Printf("Ignored post \"%s\": %s\n", watcher.name, news.Title)
			continue
		}
		log.Printf("New post from \"%s\": %s\n", watcher.name, news.Title)
		ch <- common.NotifyPayload{
			Title: watcher.name,
			Body:  news.Title,
			URL:   fmt.Sprintf("%s%s", "https://ak.hypergryph.com/news/", news.Cid),
		}
	}
}
