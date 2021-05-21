package watcher

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly/v2"
	"github.com/hguandl/dr-feeder/v2/common"
	"github.com/mitchellh/mapstructure"
)

const safariUA = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6)" +
	" AppleWebKit/605.1.15 (KHTML, like Gecko)" +
	" Version/14.0.3 Safari/605.1.15"

const iOSWeiboIntlUA = "WeiboOverseas/4.2.6 (iPhone; iOS 14.4.2; Scale/3.00)"

type weiboWatcher struct {
	uid         uint64
	updateTime  time.Time
	containerID string
	name        string
	latestMblog mblog
	sub         string
	debugURL    string
}

// NewWeiboWatcher creates a Watcher of Arknights official Weibo.
func NewWeiboWatcher(uid int64, sub string, debugURL string) (Watcher, error) {
	watcher := new(weiboWatcher)
	watcher.uid = uint64(uid)
	watcher.sub = sub
	watcher.updateTime = time.Now().UTC()
	watcher.debugURL = debugURL
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
	if watcher.debugURL != "" {
		return watcher.debugURL
	}

	return fmt.Sprintf("%s%s%d%s%s",
		"https://m.weibo.cn/api/container/getIndex?type=uid",
		"&value=", watcher.uid,
		"&containerid=", watcher.containerID,
	)
}

func (watcher weiboWatcher) weiboIntlShareAPI(weiboID string) string {
	return fmt.Sprintf("https://weibointl.api.weibo.cn/portal.php"+
		"?a=get_share_url"+
		"&ct=weibo"+
		"&lang=zh-Hans"+
		"&uid=%d"+
		"&weibo_id=%s",
		watcher.uid,
		weiboID,
	)
}

func (watcher weiboWatcher) weiboIntlShareURL(weiboID string) (string, error) {
	var err error = nil
	var shareURL string
	c := colly.NewCollector(
		colly.UserAgent(iOSWeiboIntlUA),
	)

	c.OnError(func(_ *colly.Response, e error) {
		err = e
	})

	c.OnResponse(func(r *colly.Response) {
		var resp weiboIntlResp
		err = json.Unmarshal(r.Body, &resp)
		if err != nil {
			return
		}
		if resp.Retcode != 0 {
			err = fmt.Errorf("%s", resp.Info)
			return
		}
		shareURL = resp.Data.URL
	})

	c.Visit(watcher.weiboIntlShareAPI(weiboID))
	c.Wait()
	return shareURL, err
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

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", fmt.Sprintf("%v=%v", "SUB", watcher.sub))
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

func (watcher weiboWatcher) parseContent() (common.NotifyPayload, bool) {
	weibo := watcher.latestMblog

	if len(weibo.RetweetedStatus) > 0 {
		return common.NotifyPayload{}, false
	}

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

	if strings.Contains(texts, "微博官方唯一抽奖工具") && strings.Contains(texts, "结果公正有效") {
		return common.NotifyPayload{}, false
	}

	picURL := weibo.PicURL
	pageURL, err := watcher.weiboIntlShareURL(weibo.ID)
	if err != nil {
		log.Println(err)
		pageURL = fmt.Sprintf("%s/%s", "https://m.weibo.cn/status", weibo.ID)
	}

	var pageInfo pageInfo
	mapstructure.Decode(weibo.PageInfo, &pageInfo)

	if pageInfo.Type == "article" || pageInfo.Type == "video" {
		picURL = pageInfo.PagePic.URL
	}

	if pageInfo.Type == "article" {
		pageURL = pageInfo.PageURL
	}

	return common.NotifyPayload{
		Title:  watcher.name,
		Body:   texts,
		URL:    pageURL,
		PicURL: picURL,
	}, true
}

func (watcher *weiboWatcher) Produce(ch chan common.NotifyPayload) {
	if watcher.update() {
		log.Printf("New post from \"%s\"...\n", watcher.name)
		payload, valid := watcher.parseContent()
		if valid {
			ch <- payload
		} else {
			log.Println("Useless message, ignored.")
		}
	} else {
		log.Printf("Waiting for post \"%s\"...\n", watcher.name)
	}
}
