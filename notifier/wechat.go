package notifier

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/hguandl/dr-feeder/v2/common"
	"github.com/hguandl/dr-feeder/v2/notifier/wxmsgapp"
)

type workWechatNotifier struct {
	client *wxmsgapp.WxAPIClient
}

func (notifier workWechatNotifier) String() string {
	return notifier.client.String()
}

type article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Picurl      string `json:"picurl"`
}

type news struct {
	Articles []article `json:"articles"`
}

type msgPayload struct {
	Touser  string `json:"touser"`
	Msgtype string `json:"msgtype"`
	Agentid string `json:"agentid"`
	News    news   `json:"news"`
}

// FromWxAPIClient creates a Notifier with an API client.
// Compatible with hguandl/rhodes-deliver
func FromWxAPIClient(client *wxmsgapp.WxAPIClient) Notifier {
	return workWechatNotifier{client: client}
}

func formatText(payload common.NotifyPayload) (string, string) {
	var title, desc string

	firstParaIdx := strings.Index(payload.Body, "\n\n")

	// Only one paragraph
	if firstParaIdx == -1 {
		// Short content.
		if len(payload.Body) < 128 {
			title = payload.Body
			desc = "点击查看原文"
			// Long content.
		} else {
			title = payload.Title
			desc = payload.Body
		}
		return title, desc
	}

	// 1st paragraph is short. which can be seen as the title.
	if firstParaIdx <= 128 {
		title = payload.Body[:firstParaIdx]
		desc = payload.Body[firstParaIdx+2:]
		return title, desc
	}

	// 1st paragraph is too long. Use the general title.
	if firstParaIdx > 128 {
		title = payload.Title
		desc = payload.Body
		return title, desc
	}

	// Default results
	return payload.Title, payload.Body
}

func (notifier workWechatNotifier) Push(payload common.NotifyPayload) {
	title, desc := formatText(payload)

	articles := make([]article, 1)
	articles[0] = article{
		Title:       title,
		Description: desc,
		URL:         payload.URL,
		Picurl:      payload.PicURL,
	}
	data, err := json.Marshal(
		msgPayload{
			Touser:  notifier.client.ToUser,
			Msgtype: "news",
			Agentid: notifier.client.AgentID,
			News:    news{Articles: articles},
		},
	)
	if err != nil {
		log.Println(err)
		return
	}

	err = notifier.client.SendMsg(data)
	if err != nil {
		log.Println(err)
	}
}
