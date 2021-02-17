package main

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/hguandl/arknights-news-watcher/v2/notifier"
	"gopkg.in/yaml.v2"
)

type yamlConfig struct {
	Version   string
	Notifiers []map[string]interface{}
}

func ParseConfig(path string) ([]notifier.Notifier, error) {
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

	ret := make([]notifier.Notifier, len(config.Notifiers))
	for idx, ntfc := range config.Notifiers {
		ntft, ok := ntfc["type"].(string)
		if !ok {
			err = errors.New("Invalid notifier config")
			break
		}

		switch ntft {
		case "custom":
			ntf, ok := notifier.FromCustomNotifierConfig(ntfc)
			if ok {
				ret[idx] = ntf
			} else {
				err = fmt.Errorf("Cannot parse notifier #%d", idx)
			}
		case "tgbot":
			ntf, ok := notifier.FromTgBotNotifierConfig(ntfc)
			if ok {
				ret[idx] = ntf
			} else {
				err = fmt.Errorf("Cannot parse notifier #%d", idx)
			}
		case "workwx":
			ntf, ok := notifier.FromWorkWechatNotifierConfig(ntfc)
			if ok {
				ret[idx] = ntf
			} else {
				err = fmt.Errorf("Cannot parse notifier #%d", idx)
			}
		case "bark":
			ntf, ok := notifier.FromBarkNotifierConfig(ntfc)
			if ok {
				ret[idx] = ntf
			} else {
				err = fmt.Errorf("Cannot parse notifier #%d", idx)
			}
		case "ifttt":
			ntf, ok := notifier.FromIFTTTNotifierConfig(ntfc)
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
