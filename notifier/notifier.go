package notifier

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/hguandl/dr-feeder/v2/common"
	"github.com/hguandl/dr-feeder/v2/notifier/wxmsgapp"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

// Notifier is a way to push messages to devices.
type Notifier interface {
	Push(common.NotifyPayload)
}

type yamlConfig struct {
	Version   string
	Notifiers []map[string]interface{}
}

// ParseConfig loads from config file and returns a list of Notifiers.
func ParseConfig(path string) ([]Notifier, error) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config yamlConfig
	err = yaml.Unmarshal([]byte(yamlFile), &config)
	if err != nil {
		return nil, err
	}

	if config.Version != "1.0" {
		return nil, errors.New("Invalid config version")
	}

	ret := make([]Notifier, len(config.Notifiers))
	for idx, ntfc := range config.Notifiers {
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
