package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hguandl/dr-feeder/v2/common"
)

type workWechatNotifier struct {
	corpID     string
	agentID    int
	corpSecret string
	toUser     string
}

type textCardPayload struct {
	Touser   string   `json:"touser"`
	Msgtype  string   `json:"msgtype"`
	Agentid  int      `json:"agentid"`
	Textcard textCard `json:"textcard"`
}

type textCard struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Btntxt      string `json:"btntxt"`
}

func NewWorkWechatNotifier(corpID string, agentID int, corpSecret string,
	toUser string) Notifier {
	return workWechatNotifier{
		corpID:     corpID,
		agentID:    agentID,
		corpSecret: corpSecret,
		toUser:     toUser,
	}
}

func FromWorkWechatNotifierConfig(config map[string]interface{}) (Notifier, bool) {
	corpID, ok := config["corpid"].(string)
	if !ok {
		return nil, false
	}

	agentID, ok := config["agentid"].(int)
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

func (notifier workWechatNotifier) apiURL() string {
	return fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send"+
		"?access_token=%s", notifier.corpSecret,
	)
}

func (notifier workWechatNotifier) Push(payload common.NotifyPayload) {
	var title, desc string
	if len(payload.Body) > 64 {
		title = payload.Body[:64]
		desc = payload.Body[64:]
	} else {
		title = payload.Body[:64]
		desc = "点击查看全文"
	}

	data, err := json.Marshal(
		textCardPayload{
			Touser:  notifier.toUser,
			Msgtype: "textcard",
			Agentid: notifier.agentID,
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

	_, err = http.Post(notifier.apiURL(),
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		log.Println(err)
	}
}
