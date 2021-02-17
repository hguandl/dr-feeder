package watcher

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly/v2"
	"github.com/hguandl/arknights-news-watcher/v2/common"
)

const safariUA = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6)" +
	" AppleWebKit/605.1.15 (KHTML, like Gecko)" +
	" Version/14.0.3 Safari/605.1.15"

type indexData struct {
	Data struct {
		UserInfo struct {
			ScreenName string `json:"screen_name"`
		} `json:"userInfo"`
		TabsInfo struct {
			Tabs []struct {
				TabType     string `json:"tab_type"`
				Containerid string `json:"containerid"`
			} `json:"tabs"`
		} `json:"tabsInfo"`
	} `json:"data"`
}

type mblog struct {
	CreatedAt string `json:"created_at"`
	ID        string `json:"id"`
	Text      string `json:"text"`
}

type cardData struct {
	Data struct {
		Cards []struct {
			CardType int   `json:"card_type"`
			Mblog    mblog `json:"mblog,omitempty"`
		} `json:"cards"`
	} `json:"data"`
}

type weiboWatcher struct {
	uid         uint64
	updateTime  time.Time
	containerID string
	name        string
	latestMblog mblog
}

func NewWeiboWatcher(uid int64) (Watcher, error) {
	watcher := new(weiboWatcher)
	watcher.uid = uint64(uid)
	watcher.updateTime = time.Now().UTC()
	err := watcher.setup()
	return watcher, err
}

func (watcher weiboWatcher) homeAPI() string {
	return fmt.Sprintf("%s%s%d",
		"https://m.weibo.cn/api/container/getIndex?type=uid",
		"&value=", watcher.uid,
	)
}

func (watcher weiboWatcher) weiboAPI() string {
	return fmt.Sprintf("%s%s%d%s%s",
		"https://m.weibo.cn/api/container/getIndex?type=uid",
		"&value=", watcher.uid,
		"&containerid=", watcher.containerID,
	)
}

func (watcher *weiboWatcher) setup() error {
	var err error = nil
	c := colly.NewCollector(
		colly.UserAgent(safariUA),
	)

	c.OnError(func(_ *colly.Response, e error) {
		err = e
	})

	c.OnResponse(func(r *colly.Response) {
		var data indexData
		err = json.Unmarshal(r.Body, &data)
		if err != nil {
			return
		}
		watcher.name = "微博 " + data.Data.UserInfo.ScreenName
		for _, tab := range data.Data.TabsInfo.Tabs {
			if tab.TabType == "weibo" {
				watcher.containerID = tab.Containerid
				break
			}
		}
	})

	c.Visit(watcher.homeAPI())
	c.Wait()
	return err
}

func (watcher *weiboWatcher) update() bool {
	var err error = nil
	ret := false
	c := colly.NewCollector(
		colly.UserAgent(safariUA),
	)

	c.OnError(func(_ *colly.Response, e error) {
		err = e
	})

	c.OnResponse(func(r *colly.Response) {
		var data cardData
		err = json.Unmarshal(r.Body, &data)
		if err != nil {
			return
		}
		for _, card := range data.Data.Cards {
			if card.CardType == 9 {
				var dateTime time.Time
				dateTime, err = time.Parse(time.RubyDate, card.Mblog.CreatedAt)
				if dateTime.After(watcher.updateTime) {
					ret = true
					watcher.updateTime = dateTime
					watcher.latestMblog = card.Mblog
				}
			}
		}
	})

	c.Visit(watcher.weiboAPI())
	c.Wait()

	if err != nil {
		log.Println(err)
		ret = false
	}
	return ret
}

func (watcher weiboWatcher) parseContent() common.NotifyPayload {
	weibo := watcher.latestMblog

	doc, _ := htmlquery.Parse(
		strings.NewReader(
			strings.ReplaceAll(weibo.Text, "<br />", "\n"),
		),
	)
	nodes, _ := htmlquery.QueryAll(doc, "//text()")

	texts := ""
	for _, node := range nodes {
		if node.Data == "#明日方舟#" {
			continue
		}
		texts += strings.Trim(node.Data, " \n")
	}

	return common.NotifyPayload{
		Title: watcher.name,
		Body:  texts,
		URL:   fmt.Sprintf("%s/%s", "https://m.weibo.cn/status", weibo.ID),
	}
}

func (watcher *weiboWatcher) Produce(ch chan common.NotifyPayload) {
	if watcher.update() {
		log.Printf("New post from \"%s\"...\n", watcher.name)
		ch <- watcher.parseContent()
	} else {
		log.Printf("Waiting for post \"%s\"...\n", watcher.name)
	}
}
