package main_test

import (
	"testing"

	ak "github.com/hguandl/dr-feeder/v2"
)

func TestParseConfig(t *testing.T) {
	notifiers, err := ak.ParseConfig("config.yaml")

	if err != nil {
		t.Error(err)
	}

	t.Log(notifiers)
}
