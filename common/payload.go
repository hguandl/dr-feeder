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

// UTF8TruncRunes returns the first <length> characters of <src> string with UTF-8 runes.
func UTF8TruncRunes(src string, length int) string {
	var runeSize, readBytes int = 0, 0

	for i := 0; i < length && readBytes < len(src); i++ {
		_, runeSize = utf8.DecodeRuneInString(src[readBytes:])
		readBytes += runeSize
	}

	return src[:readBytes]
}

// UTF8TruncBytesByRunes returns the first <size> bytes of <src> string with UTF-8 runes.
func UTF8TruncBytesByRunes(src string, size int) string {
	var runeSize, readBytes int = 0, 0

	if size > len(src) {
		size = len(src)
	}

	for readBytes < size {
		_, runeSize = utf8.DecodeRuneInString(src[readBytes:])
		readBytes += runeSize
	}
	return src[:readBytes]
}

func (payload NotifyPayload) String() string {
	body := UTF8TruncRunes(payload.Body, 20)
	body = strings.ReplaceAll(body, "\n", "")

	return fmt.Sprintf("{Title: %v} {Body: %v} {URL: %v} {PicURL: %v}",
		payload.Title,
		body,
		payload.URL,
		payload.PicURL,
	)
}
