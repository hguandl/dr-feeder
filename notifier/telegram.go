package notifier

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/hguandl/dr-feeder/v2/common"
)

type tgBotNotifier struct {
	BotAPIKey string `mapstructure:"api_key"`
	Chats     []string
}

func (notifier tgBotNotifier) apiURL() string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage",
		notifier.BotAPIKey,
	)
}

func (notifier tgBotNotifier) Push(payload common.NotifyPayload) {
	texts := payload.Body + "\n\n" + payload.URL

	for _, chat := range notifier.Chats {
		r, err := http.PostForm(notifier.apiURL(),
			url.Values{
				"chat_id": {fmt.Sprint(chat)},
				"text":    {texts},
			})
		if err != nil {
			log.Println(err)
		} else {
			r.Body.Close()
		}
	}
}
