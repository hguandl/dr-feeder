package watcher

import "github.com/hguandl/dr-feeder/v2/common"

type Watcher interface {
	Produce(chan common.NotifyPayload)
}
