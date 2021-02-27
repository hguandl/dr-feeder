package notifier

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/hguandl/dr-feeder/v2/common"
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
			ntf, ok := FromCustomNotifierConfig(ntfc)
			if ok {
				ret[idx] = ntf
			} else {
				err = fmt.Errorf("Cannot parse notifier #%d", idx)
			}
		case "tgbot":
			ntf, ok := FromTgBotNotifierConfig(ntfc)
			if ok {
				ret[idx] = ntf
			} else {
				err = fmt.Errorf("Cannot parse notifier #%d", idx)
			}
		case "workwx":
			ntf, ok := FromWorkWechatNotifierConfig(ntfc)
			if ok {
				ret[idx] = ntf
			} else {
				err = fmt.Errorf("Cannot parse notifier #%d", idx)
			}
		case "bark":
			ntf, ok := FromBarkNotifierConfig(ntfc)
			if ok {
				ret[idx] = ntf
			} else {
				err = fmt.Errorf("Cannot parse notifier #%d", idx)
			}
		case "ifttt":
			ntf, ok := FromIFTTTNotifierConfig(ntfc)
			if ok {
				ret[idx] = ntf
			} else {
				err = fmt.Errorf("Cannot parse notifier #%d", idx)
			}
		default:
			err = fmt.Errorf("Unknown notifier type \"%s\"", ntft)
		}

		if err != nil {
			break
		}
	}

	return ret, err
}
