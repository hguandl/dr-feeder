package watcher

import "github.com/hguandl/dr-feeder/v2/common"

// Watcher is a news source to watch.
type Watcher interface {
	Produce(chan common.NotifyPayload)
}
