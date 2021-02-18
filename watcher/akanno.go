package watcher

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/hguandl/dr-feeder/v2/common"
)

const iOSClientUA = "arknights/385" +
	" CFNetwork/1220.1" +
	" Darwin/20.3.0'"

type announce struct {
	AnnounceID string `json:"announceId"`
	Title      string `json:"title"`
	IsWebURL   bool   `json:"isWebUrl"`
	WebURL     string `json:"webUrl"`
	Day        int    `json:"day"`
	Month      int    `json:"month"`
	Group      string `json:"group"`
}

type announceMeta struct {
	FocusAnnounceID string     `json:"focusAnnounceId"`
	AnnounceList    []announce `json:"announceList"`
}

type akAnnounceWatcher struct {
	name       string
	latestID   string
	latestAnno announce
}

func NewAkAnnounceWatcher() (Watcher, error) {
	watcher := new(akAnnounceWatcher)
	watcher.name = "明日方舟客户端公告"
	err := watcher.setup()
	return watcher, err
}

func (watcher akAnnounceWatcher) fetchAPI() (announceMeta, error) {
	const apiURL = "https://ak-fs.hypergryph.com/announce/IOS/announcement.meta.json?sign="
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

	c.Visit(apiURL)
	c.Wait()

	return data, err
}

func (watcher *akAnnounceWatcher) setup() error {
	data, err := watcher.fetchAPI()
	if err != nil {
		return err
	}

	watcher.latestID = data.FocusAnnounceID

	return nil
}

func (watcher *akAnnounceWatcher) update() bool {
	data, err := watcher.fetchAPI()
	if err != nil {
		log.Println(err)
		return false
	}

	if data.FocusAnnounceID != watcher.latestID {
		watcher.latestID = data.FocusAnnounceID
		for _, anno := range data.AnnounceList {
			if anno.AnnounceID == data.FocusAnnounceID {
				watcher.latestAnno = anno
				break
			}
		}
		if strings.Contains(watcher.latestAnno.Title, "制作组通讯") {
			return true
		}
	}

	return false
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
