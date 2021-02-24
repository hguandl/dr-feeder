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

type textCardPayload struct {
	Touser   string   `json:"touser"`
	Msgtype  string   `json:"msgtype"`
	Agentid  string   `json:"agentid"`
	Textcard textCard `json:"textcard"`
}

type textCard struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Btntxt      string `json:"btntxt"`
}

// NewWorkWechatNotifier creates a Notifier with Work Wechat App.
func NewWorkWechatNotifier(corpID string, agentID string, corpSecret string,
	toUser string) Notifier {
	client := wxmsgapp.WxAPIClient{
		CorpID:     corpID,
		ToUser:     toUser,
		AgentID:    agentID,
		CorpSecret: corpSecret,
	}
	return workWechatNotifier{client: &client}
}

// FromWorkWechatNotifierConfig parses the config to create a workWechatNotifier.
func FromWorkWechatNotifierConfig(config map[string]interface{}) (Notifier, bool) {
	corpID, ok := config["corpid"].(string)
	if !ok {
		return nil, false
	}

	agentID, ok := config["agentid"].(string)
	if !ok {
		return nil, false
	}

	corpSecret, ok := config["corpsecret"].(string)
	if !ok {
		return nil, false
	}

	toUser, ok := config["touser"].(string)
	if !ok {
		return nil, false
	}

	return NewWorkWechatNotifier(corpID, agentID, corpSecret, toUser), true
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

	data, err := json.Marshal(
		textCardPayload{
			Touser:  notifier.client.ToUser,
			Msgtype: "textcard",
			Agentid: notifier.client.AgentID,
			Textcard: textCard{
				Title:       title,
				Description: desc,
				URL:         payload.URL,
				Btntxt:      "全文",
			},
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
