package notifier

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hguandl/arknights-news-watcher/v2/common"
)

type barkNotifier struct {
	apiURL    string
	apiTokens []string
}

func NewBarkNotifier(apiTokens []string) Notifier {
	return barkNotifier{
		apiURL:    "https://api.day.app",
		apiTokens: apiTokens,
	}
}

func FromBarkNotifierConfig(config map[string]interface{}) (Notifier, bool) {
	tokensRaw, ok := config["tokens"].([]interface{})
	if !ok {
		return nil, false
	}

	apiTokens := make([]string, len(tokensRaw))
	for i := range apiTokens {
		apiTokens[i], ok = tokensRaw[i].(string)
		if !ok {
			return nil, false
		}
	}

	return NewBarkNotifier(apiTokens), true
}

func (notifier barkNotifier) Push(payload common.NotifyPayload) {
	for _, token := range notifier.apiTokens {
		_, err := http.Get(fmt.Sprintf(
			"%s/%s/%s/%s?url=%s",
			notifier.apiURL,
			token,
			payload.Title,
			payload.Body,
			payload.URL,
		))
		if err != nil {
			log.Println(err)
		}
	}
}
