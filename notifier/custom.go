package notifier

import (
	"log"
	"net/http"
	"net/url"

	"github.com/hguandl/arknights-news-watcher/v2/common"
)

type customNotifer struct {
	apiURL string
}

func NewCustomNotifier(apiURL string) Notifier {
	var notifer customNotifer
	notifer.apiURL = apiURL
	return notifer
}

func FromCustomNotifierConfig(config map[string]interface{}) (Notifier, bool) {
	apiURL, ok := config["api_url"].(string)
	if !ok {
		return nil, false
	}

	return NewCustomNotifier(apiURL), true
}

func (notifer customNotifer) Push(payload common.NotifyPayload) {
	_, err := http.PostForm(notifer.apiURL,
		url.Values{
			"title": {payload.Title},
			"body":  {payload.Body},
			"url":   {payload.URL},
		})
	if err != nil {
		log.Println(err)
	}
}
