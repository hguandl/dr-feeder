package notifier

import (
	"log"
	"net/http"
	"net/url"

	"github.com/hguandl/dr-feeder/v2/common"
)

type customNotifier struct {
	APIURL string `mapstructure:"api_url"`
}

func (notifier customNotifier) Push(payload common.NotifyPayload) {
	r, err := http.PostForm(notifier.APIURL,
		url.Values{
			"title":  {payload.Title},
			"body":   {payload.Body},
			"url":    {payload.URL},
			"picurl": {payload.PicURL},
		})
	if err != nil {
		log.Println(err)
		return
	}
	defer r.Body.Close()
}
