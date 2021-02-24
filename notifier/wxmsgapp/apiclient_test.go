package wxmsgapp_test

import (
	"encoding/json"
	"testing"

	"github.com/hguandl/dr-feeder/v2/notifier/wxmsgapp"
)

func TestJSONUnmarshal(t *testing.T) {
	const testStr = `{
	"agentid": "100002",
	"touser": "@all",
	"corpid": "rhodesisland",
	"CorpSecret": "foobar******"
}`

	var data wxmsgapp.WxAPIClient
	err := json.Unmarshal([]byte(testStr), &data)
	if err != nil {
		t.Error(err)
	}

	t.Log(data.AccessToken)
	t.Log(data.TokenUntil)
}
