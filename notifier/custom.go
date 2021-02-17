package notifier

import (
	"log"
	"net/http"
	"net/url"

	"github.com/hguandl/arknights-news-watcher/v2/common"
)

type customNotifier struct {
	apiURL string
}

func NewCustomNotifier(apiURL string) Notifier {
	return customNotifier{
		apiURL: apiURL,
	}
}

func FromCustomNotifierConfig(config map[string]interface{}) (Notifier, bool) {
	apiURL, ok := config["api_url"].(string)
	if !ok {
		return nil, false
	}

	return NewCustomNotifier(apiURL), true
}

func (notifier customNotifier) Push(payload common.NotifyPayload) {
	_, err := http.PostForm(notifier.apiURL,
		url.Values{
			"title": {payload.Title},
			"body":  {payload.Body[:128]},
			"url":   {payload.URL},
		})
	if err != nil {
		log.Println(err)
	}
}
