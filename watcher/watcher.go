package watcher

import "github.com/hguandl/arknights-news-watcher/v2/common"

type Watcher interface {
	Produce(chan common.NotifyPayload)
}
