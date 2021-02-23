package common

// NotifyPayload contains contents of a message to push.
// Title is the source of the news.
// Body is the summary of the message.
// URL is the link to detailed information.
type NotifyPayload struct {
	Title string
	Body  string
	URL   string
}
