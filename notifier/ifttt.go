package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hguandl/dr-feeder/v2/common"
)

type iftttNotifier struct {
	Webhooks []webhookConfig
}

type webhookConfig struct {
	Event  string
	APIKey string `mapstructure:"api_key"`
}

type webhookPayload struct {
	Value1 string `json:"value1"`
	Value2 string `json:"value2"`
}

func (notifier iftttNotifier) apiURL(webhook webhookConfig) string {
	return fmt.Sprintf("https://make.ifttt.com/trigger/%s/with/key/%s",
		webhook.Event,
		webhook.APIKey)
}

func (notifier iftttNotifier) Push(payload common.NotifyPayload) {
	for _, webhook := range notifier.Webhooks {
		data, err := json.Marshal(
			webhookPayload{
				Value1: payload.Body,
				Value2: payload.URL,
			},
		)
		if err != nil {
			log.Println(err)
			return
		}

		r, err := http.Post(
			notifier.apiURL(webhook),
			"application/json",
			bytes.NewBuffer(data),
		)
		if err != nil {
			log.Println(err)
		} else {
			r.Body.Close()
		}
	}
}
