package notifier

import "github.com/hguandl/arknights-news-watcher/v2/common"

type Notifier interface {
	Push(common.NotifyPayload)
}
