package main_test

import (
	"testing"

	"github.com/hguandl/dr-feeder/v2/notifier"
)

func TestParseConfig(t *testing.T) {
	notifiers, err := notifier.ParseConfig("config.yaml")

	if err != nil {
		t.Error(err)
	}

	for _, n := range notifiers {
		t.Logf("%v\n", n)
	}
}
