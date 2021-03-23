package main_test

import (
	"testing"

	main "github.com/hguandl/dr-feeder/v2"
	"github.com/hguandl/dr-feeder/v2/notifier"
	"github.com/hguandl/dr-feeder/v2/watcher"
)

var config main.YamlConfig

func TestLoadConfig(t *testing.T) {
	var err error = nil

	config, err = main.LoadConfig("config.yaml")
	if err != nil {
		t.Error(err)
	}

	t.Logf("Config file version %v.", config.Version)
}

func TestParseNotifiers(t *testing.T) {
	notifiers, err := notifier.ParseNotifiers(config.Notifiers)

	if err != nil {
		t.Error(err)
	}

	for _, n := range notifiers {
		t.Logf("%v", n)
	}
}

func TestParseWatchers(t *testing.T) {
	watchers, err := watcher.ParseWatchers(config.Watchers, "", false)

	if err != nil {
		t.Error(err)
	}

	for _, n := range watchers {
		t.Logf("%v", n)
	}
}
