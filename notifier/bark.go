package notifier

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hguandl/dr-feeder/v2/common"
)

type barkNotifier struct {
	apiURL    string
	apiTokens []string
}

// NewBarkNotifier creates a Notifier with iOS Bark App.
func NewBarkNotifier(apiTokens []string) Notifier {
	return barkNotifier{
		apiURL:    "https://api.day.app",
		apiTokens: apiTokens,
	}
}

// FromBarkNotifierConfig parses the config file to create a barkNotifier.
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
		r, err := http.Get(fmt.Sprintf(
			"%s/%s/%s/%s?url=%s",
			notifier.apiURL,
			token,
			payload.Title,
			payload.Body,
			payload.URL,
		))
		if err != nil {
			log.Println(err)
		} else {
			r.Body.Close()
		}
	}
}
