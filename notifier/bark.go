package notifier

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/hguandl/dr-feeder/v2/common"
)

type barkNotifier struct {
	Tokens []string
}

type barkPayload struct {
	Title     string `json:"title"`
	Body      string `json:"body"`
	Category  string `json:"category"`
	DeviceKey string `json:"device_key"`
	URL       string `json:"url"`
}

func (notifier barkNotifier) apiURL() string {
	return "https://api.day.app/push"
}

func (notifier barkNotifier) Push(payload common.NotifyPayload) {
	for _, token := range notifier.Tokens {
		pushPayload := barkPayload{
			Title:     payload.Title,
			Body:      payload.Body,
			Category:  "",
			DeviceKey: token,
			URL:       payload.URL,
		}

		postPayload, err := json.Marshal(pushPayload)
		if err != nil {
			log.Println("JSON: ", err)
			continue
		}
		postBody := bytes.NewBuffer(postPayload)

		r, err := http.Post(notifier.apiURL(), "application/json; charset=utf-8", postBody)
		if err != nil {
			log.Println("POST: ", err)
		} else {
			r.Body.Close()
		}
	}
}
