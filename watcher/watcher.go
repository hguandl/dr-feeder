package watcher

import (
	"errors"
	"fmt"
	"path"

	"github.com/hguandl/dr-feeder/v2/common"
	"github.com/mitchellh/mapstructure"
)

// Watcher is a news source to watch.
type Watcher interface {
	Produce(chan common.NotifyPayload)
}

type weiboConfig struct {
	UID      int64
	Sub      string
	DebugURL string `mapstructure:"debug_url"`
}

type akAnnoConfig struct {
	Channel  string
	DebugURL string `mapstructure:"debug_url"`
}

type sirenConfig struct {
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
func ParseWatchers(configs []map[string]interface{}, dataPath string, debugMode bool) ([]Watcher, error) {
	var err error = nil
	ret := make([]Watcher, len(configs))

	for idx, config := range configs {
		watcherType, ok := config["type"].(string)
		if !ok {
			err = errors.New("invalid watcher config")
			break
		}

		switch watcherType {
		case "weibo":
			var wbConfig weiboConfig
			if err = mapstructure.Decode(config, &wbConfig); err != nil {
				break
			}
			ret[idx], err = NewWeiboWatcher(wbConfig.UID, wbConfig.Sub, wrapDebug(wbConfig.DebugURL, debugMode))
		case "akanno":
			var akConfig akAnnoConfig
			if err = mapstructure.Decode(config, &akConfig); err != nil {
				break
			}
			if akConfig.Channel != "IOS" {
				err = fmt.Errorf("unsupported channel \"%v\"", akConfig.Channel)
				break
			}
			ret[idx], err = NewAkAnnounceWatcher(path.Join(dataPath, "akanno.db"), wrapDebug(akConfig.DebugURL, debugMode))
		case "siren":
			var akConfig sirenConfig
			if err = mapstructure.Decode(config, &akConfig); err != nil {
				break
			}
			ret[idx], err = NewSirenWatcher(path.Join(dataPath, "siren.db"), wrapDebug(akConfig.DebugURL, debugMode))
		default:
			err = fmt.Errorf("unknown watcher #%d with type \"%s\"", idx, watcherType)
		}

		if err != nil {
			break
		}
	}

	return ret, err
}
