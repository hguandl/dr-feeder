package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hguandl/arknights-news-watcher/v2/common"
)

type iftttNotifier struct {
	webhooks []webhookConfig
}

type webhookConfig struct {
	event  string
	apiKey string
}

type webhookPayload struct {
	Value1 string `json:"value1"`
	Value2 string `json:"value2"`
}

func NewIFTTTNotifier(webhooks []webhookConfig) Notifier {
	return iftttNotifier{
		webhooks: webhooks,
	}
}

func FromIFTTTNotifierConfig(config map[string]interface{}) (Notifier, bool) {
	listRaw, ok := config["webhooks"].([]interface{})
	if !ok {
		return nil, false
	}

	webhooks := make([]webhookConfig, len(listRaw))
	for i := range webhooks {
		whRaw, ok := listRaw[i].(map[interface{}]interface{})
		if !ok {
			return nil, false
		}

		var wh webhookConfig
		for k, v := range whRaw {
			strK, ok := k.(string)
			strV, ok := v.(string)
			if !ok {
				return nil, false
			}

			switch strK {
			case "event":
				wh.event = strV
			case "api_key":
				wh.apiKey = strV
			}
		}

		webhooks[i] = wh
	}

	return NewIFTTTNotifier(webhooks), true
}

func (notifier iftttNotifier) apiURL(webhook webhookConfig) string {
	return fmt.Sprintf("https://make.ifttt.com/trigger/%s/with/key/%s",
		webhook.event,
		webhook.apiKey)
}

func (notifier iftttNotifier) Push(payload common.NotifyPayload) {
	for _, webhook := range notifier.webhooks {

		data, err := json.Marshal(
			webhookPayload{
				Value1: payload.Body,
				Value2: payload.URL,
			},
		)

		_, err = http.Post(
			notifier.apiURL(webhook),
			"application/json",
			bytes.NewBuffer(data),
		)
		if err != nil {
			log.Println(err)
		}
	}
}
