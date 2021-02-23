package notifier

import "github.com/hguandl/dr-feeder/v2/common"

// Notifier is a way to push messages to devices.
type Notifier interface {
	Push(common.NotifyPayload)
}
