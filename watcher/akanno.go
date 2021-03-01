package watcher

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/hguandl/dr-feeder/v2/common"
)

const iOSClientUA = "arknights/385" +
	" CFNetwork/1220.1" +
	" Darwin/20.3.0'"

type akAnnounceWatcher struct {
	name       string
	focusID    string
	latestAnno announce
	existedID  []string
	debugURL   string
}

// NewAkAnnounceWatcher creates a Watcher of Arknights game annoucements.
func NewAkAnnounceWatcher(debugURL string) (Watcher, error) {
	watcher := new(akAnnounceWatcher)
	watcher.name = "明日方舟客户端公告"
	watcher.debugURL = debugURL
	err := watcher.setup()
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

	watcher.focusID = data.FocusAnnounceID
	watcher.existedID = flushIDList(data.AnnounceList)

	return nil
}

func flushIDList(announceList []announce) []string {
	ret := make([]string, len(announceList))
	for i, anno := range announceList {
		ret[i] = anno.AnnounceID
	}

	return ret
}

func (watcher *akAnnounceWatcher) update() bool {
	data, err := watcher.fetchAPI()
	if err != nil {
		log.Println(err)
		return false
	}

	if watcher.focusID != data.FocusAnnounceID {
		watcher.focusID = data.FocusAnnounceID
		existed := false
		for _, anno := range data.AnnounceList {
			if anno.AnnounceID == data.FocusAnnounceID {
				existed = true
				break
			}
		}
		if existed == false {
			watcher.latestAnno = announce{
				Title:  "出现公告弹窗，可能会有新饼",
				WebURL: "https://ak.hypergryph.com/news.html",
			}
			return true
		}
	}

	for _, anno := range data.AnnounceList {
		newID := anno.AnnounceID
		existed := false
		for _, oldID := range watcher.existedID {
			if newID == oldID {
				existed = true
				break
			}
		}
		if existed == false {
			watcher.existedID = flushIDList(data.AnnounceList)
			if strings.Contains(anno.Title, "制作组通讯") {
				watcher.latestAnno = anno
				return true
			}
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
