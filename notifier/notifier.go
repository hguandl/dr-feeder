package notifier

import (
	"errors"
	"fmt"

	"github.com/hguandl/dr-feeder/v2/common"
	"github.com/hguandl/dr-feeder/v2/notifier/wxmsgapp"
	"github.com/mitchellh/mapstructure"
)

// Notifier is a way to push messages to devices.
type Notifier interface {
	Push(common.NotifyPayload)
}

// ParseNotifiers decodes the config returns a list of Notifiers.
func ParseNotifiers(configs []map[string]interface{}) ([]Notifier, error) {
	var err error = nil
	ret := make([]Notifier, len(configs))

	for idx, ntfc := range configs {
		ntft, ok := ntfc["type"].(string)
		if !ok {
			err = errors.New("Invalid notifier config")
			break
		}

		switch ntft {
		case "custom":
			var ntf customNotifier
			err = mapstructure.Decode(ntfc, &ntf)
			ret[idx] = ntf
		case "tgbot":
			var ntf tgBotNotifier
			err = mapstructure.Decode(ntfc, &ntf)
			ret[idx] = ntf
		case "workwx":
			var wxConfig wxmsgapp.WxAPIClient
			err = mapstructure.Decode(ntfc, &wxConfig)
			ret[idx] = workWechatNotifier{client: &wxConfig}
		case "bark":
			var ntf barkNotifier
			err = mapstructure.Decode(ntfc, &ntf)
			ret[idx] = ntf
		case "ifttt":
			var ntf iftttNotifier
			err = mapstructure.Decode(ntfc, &ntf)
			ret[idx] = ntf
		default:
			err = fmt.Errorf("Unknown notifier #%d with type \"%s\"", idx, ntft)
		}

		if err != nil {
			break
		}
	}

	return ret, err
}
