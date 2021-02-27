package watcher

import (
	"errors"
	"fmt"

	"github.com/hguandl/dr-feeder/v2/common"
	"github.com/mitchellh/mapstructure"
)

// Watcher is a news source to watch.
type Watcher interface {
	Produce(chan common.NotifyPayload)
}

type weiboConfig struct {
	UID      int64
	DebugURL string `mapstructure:"debug_url"`
}

type akAnnoConfig struct {
	Channel  string
	DebugURL string `mapstructure:"debug_url"`
}

func wrapDebug(debugURL string, debugMode bool) string {
	if debugURL != "" {
		if debugMode {
			return debugURL
		}

		println("Not on debug mode. Ignored debug URL.")
		return ""
	}
	return ""
}

// ParseWatchers decodes the config returns a list of Watchers.
func ParseWatchers(configs []map[string]interface{}, debugMode bool) ([]Watcher, error) {
	var err error = nil
	ret := make([]Watcher, len(configs))

	for idx, config := range configs {
		watcherType, ok := config["type"].(string)
		if !ok {
			err = errors.New("Invalid watcher config")
			break
		}

		switch watcherType {
		case "weibo":
			var wbConfig weiboConfig
			err = mapstructure.Decode(config, &wbConfig)
			ret[idx], err = NewWeiboWatcher(wbConfig.UID, wrapDebug(wbConfig.DebugURL, debugMode))
		case "akanno":
			var akConfig akAnnoConfig
			err = mapstructure.Decode(config, &akConfig)
			if akConfig.Channel != "IOS" {
				err = fmt.Errorf("Unsupported channel \"%v\"", akConfig.Channel)
			}
			ret[idx], err = NewAkAnnounceWatcher(wrapDebug(akConfig.DebugURL, debugMode))
		default:
			err = fmt.Errorf("Unknown watcher #%d with type \"%s\"", idx, watcherType)
		}

		if err != nil {
			break
		}
	}

	return ret, err
}
