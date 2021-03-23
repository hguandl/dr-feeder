package watcher

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/hguandl/dr-feeder/v2/common"

	bolt "go.etcd.io/bbolt"
)

const iOSClientUA = "arknights/385" +
	" CFNetwork/1220.1" +
	" Darwin/20.3.0'"

type akAnnounceWatcher struct {
	name       string
	latestAnno announce
	debugURL   string
	db         *bolt.DB
}

// NewAkAnnounceWatcher creates a Watcher of Arknights game annoucements.
func NewAkAnnounceWatcher(dbPath string, debugURL string) (Watcher, error) {
	var err error = nil

	watcher := new(akAnnounceWatcher)
	watcher.name = "明日方舟客户端公告"
	watcher.debugURL = debugURL

	watcher.db, err = bolt.Open(dbPath, 0666, nil)
	if err != nil {
		return watcher, err
	}

	err = watcher.setup()
	return watcher, err
}

func (watcher akAnnounceWatcher) apiURL() string {
	if watcher.debugURL != "" {
		return watcher.debugURL
	}
	clientID := rand.Intn(114514191) + 11451419
	return fmt.Sprintf("%s?sign=%d",
		"https://ak-fs.hypergryph.com/announce/IOS/announcement.meta.json",
		clientID,
	)
}

func (watcher akAnnounceWatcher) fetchAPI() (announceMeta, error) {
	var err error = nil
	var data announceMeta

	c := colly.NewCollector(
		colly.UserAgent(iOSClientUA),
	)

	c.OnError(func(_ *colly.Response, e error) {
		err = e
	})

	c.OnResponse(func(r *colly.Response) {
		err = json.Unmarshal(r.Body, &data)
	})

	c.Visit(watcher.apiURL())
	c.Wait()

	return data, err
}

func (watcher *akAnnounceWatcher) setup() error {
	data, err := watcher.fetchAPI()
	if err != nil {
		return err
	}

	watcher.storeAnnos(data.AnnounceList)

	return nil
}

func (watcher *akAnnounceWatcher) storeAnnos(announceList []announce) error {
	err := watcher.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("AkAnno"))
		for _, anno := range announceList {
			err = b.Put([]byte(anno.AnnounceID), []byte(anno.Title))
		}
		return err
	})

	return err
}

func (watcher *akAnnounceWatcher) update() bool {
	data, err := watcher.fetchAPI()
	if err != nil {
		log.Println(err)
		return false
	}

	ret := false
	err = watcher.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("AkAnno"))
		for _, anno := range data.AnnounceList {
			v := b.Get([]byte(anno.AnnounceID))
			if v == nil {
				if strings.Contains(anno.Title, "制作组通讯") {
					watcher.latestAnno = anno
					ret = true
				}
				err = b.Put([]byte(anno.AnnounceID), []byte(anno.Title))
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

func (watcher akAnnounceWatcher) parseContent() common.NotifyPayload {
	anno := watcher.latestAnno

	return common.NotifyPayload{
		Title: watcher.name,
		Body:  anno.Title,
		URL:   anno.WebURL,
	}
}

func (watcher *akAnnounceWatcher) Produce(ch chan common.NotifyPayload) {
	if watcher.update() {
		log.Printf("New post from \"%s\"...\n", watcher.name)
		ch <- watcher.parseContent()
	} else {
		log.Printf("Waiting for post \"%s\"...\n", watcher.name)
	}
}
