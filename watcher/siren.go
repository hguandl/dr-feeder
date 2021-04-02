package watcher

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly/v2"
	"github.com/hguandl/dr-feeder/v2/common"
	"github.com/mitchellh/mapstructure"

	bolt "go.etcd.io/bbolt"
)

type sirenWatcher struct {
	name       string
	latestNews sirenNewsData
	debugURL   string
	db         *bolt.DB
}

// NewSirenWatcher creates a Watcher of news from Monster Siren.
func NewSirenWatcher(dbPath string, debugURL string) (Watcher, error) {
	var err error = nil

	watcher := new(sirenWatcher)
	watcher.name = "塞壬唱片"
	watcher.debugURL = debugURL

	watcher.db, err = bolt.Open(dbPath, 0666, nil)
	if err != nil {
		return watcher, err
	}

	err = watcher.setup()
	return watcher, err
}

func (watcher sirenWatcher) apiURL(newsID string) string {
	if watcher.debugURL != "" {
		return watcher.debugURL
	}
	return fmt.Sprintf("%s%s",
		"https://monster-siren.hypergryph.com/api/news/",
		newsID,
	)
}

func (watcher sirenWatcher) fetchAPI(newID string) (sirenAPIPayload, error) {
	var err error = nil
	var data sirenAPIPayload

	c := colly.NewCollector(
		colly.UserAgent(safariUA),
	)

	c.OnError(func(_ *colly.Response, e error) {
		err = e
	})

	c.OnResponse(func(r *colly.Response) {
		err = json.Unmarshal(r.Body, &data)
	})

	c.Visit(watcher.apiURL(newID))
	c.Wait()

	return data, err
}

func (watcher sirenWatcher) getNewsList() ([]sirenNewsData, error) {
	data, err := watcher.fetchAPI("")
	if err != nil {
		return nil, err
	}

	var content sirenListData
	err = mapstructure.Decode(data.Data, &content)
	if err != nil {
		return nil, err
	}

	return content.List, nil
}

func (watcher sirenWatcher) getNews(newsID string) (sirenNewsData, error) {
	var content sirenNewsData
	data, err := watcher.fetchAPI(newsID)
	if err != nil {
		return content, err
	}

	err = mapstructure.Decode(data.Data, &content)
	if err != nil {
		return content, err
	}

	return content, nil
}

func (watcher *sirenWatcher) setup() error {
	newsList, err := watcher.getNewsList()
	if err != nil {
		return err
	}

	watcher.storeNews(newsList)

	return nil
}

func (watcher *sirenWatcher) storeNews(newsList []sirenNewsData) error {
	err := watcher.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("Siren"))
		for _, news := range newsList {
			err = b.Put([]byte(news.Cid), []byte(news.Title))
		}
		return err
	})

	return err
}

func (watcher *sirenWatcher) update() bool {
	newsList, err := watcher.getNewsList()
	if err != nil {
		log.Println(err)
		return false
	}

	ret := false
	err = watcher.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("Siren"))
		for _, news := range newsList {
			v := b.Get([]byte(news.Cid))
			if v == nil {
				watcher.latestNews = news
				ret = true
				err = b.Put([]byte(news.Cid), []byte(news.Title))
				break
			}
		}
		return err
	})

	if err != nil {
		log.Println(err)
		return false
	}

	return ret
}

func (watcher sirenWatcher) parseContent() common.NotifyPayload {
	news := watcher.latestNews
	fullNews, err := watcher.getNews(news.Cid)
	texts := news.Title + "\n"
	if err == nil {
		doc, _ := htmlquery.Parse(
			strings.NewReader(fullNews.Content),
		)
		nodes, _ := htmlquery.QueryAll(doc, "//text()")

		for _, node := range nodes {
			texts += "\n"
			texts += strings.Trim(node.Data, " \n")
		}
	}

	return common.NotifyPayload{
		Title: watcher.name,
		Body:  texts,
		URL:   fmt.Sprintf("%s%s", "https://monster-siren.hypergryph.com/info/", news.Cid),
	}
}

func (watcher *sirenWatcher) Produce(ch chan common.NotifyPayload) {
	if watcher.update() {
		log.Printf("New post from \"%s\"...\n", watcher.name)
		ch <- watcher.parseContent()
	} else {
		log.Printf("Waiting for post \"%s\"...\n", watcher.name)
	}
}
