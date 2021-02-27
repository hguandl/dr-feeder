package common

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// NotifyPayload contains contents of a message to push.
// Title is the source of the news.
// Body is the summary of the message.
// URL is the link to detailed information.
// PicURL is the link to the caption picture (optional).
type NotifyPayload struct {
	Title  string
	Body   string
	URL    string
	PicURL string
}

func (payload NotifyPayload) String() string {
	var size, n int = 0, 0
	for i := 0; i < 20 && n < len(payload.Body); i++ {
		_, size = utf8.DecodeRuneInString(payload.Body[n:])
		n += size
	}
	body := strings.ReplaceAll(payload.Body[:n], "\n", "")

	return fmt.Sprintf("{Title: %v} {Body: %v} {URL: %v} {PicURL: %v}",
		payload.Title,
		body,
		payload.URL,
		payload.PicURL,
	)
}
