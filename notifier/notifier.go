package notifier

import "github.com/hguandl/dr-feeder/v2/common"

type Notifier interface {
	Push(common.NotifyPayload)
}
