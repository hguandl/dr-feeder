package wxmsgapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// UpdateToken requests for an access token for current app.
func (client *WxAPIClient) UpdateToken() error {
	req, err := http.Get(
		fmt.Sprintf("%s?corpid=%s&corpsecret=%s",
			"https://qyapi.weixin.qq.com/cgi-bin/gettoken",
			client.CorpID, client.CorpSecret,
		),
	)
	if err != nil {
		return err
	}
	defer req.Body.Close()

	var resp wxAPIResp
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&resp)
	if err != nil {
		return err
	}

	if resp.Errcode != 0 {
		return fmt.Errorf(resp.Errmsg)
	}

	client.AccessToken = resp.AccessToken
	client.TokenUntil = time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second)
	return nil
}

func (client WxAPIClient) isTokenValid() bool {
	return time.Now().Before(client.TokenUntil)
}

// SendMsg sends JSON data to the Work Wechat message API.
func (client *WxAPIClient) SendMsg(msgData []byte) error {
	var err error = nil
	if !client.isTokenValid() {
		err = client.UpdateToken()
		if err != nil {
			return err
		}
	}

	req, err := http.Post(
		fmt.Sprintf("%s?access_token=%s",
			"https://qyapi.weixin.qq.com/cgi-bin/message/send",
			client.AccessToken,
		),
		"application/json",
		bytes.NewBuffer(msgData),
	)
	if err != nil {
		return err
	}
	defer req.Body.Close()

	var resp wxAPIResp
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&resp)
	if err != nil {
		return err
	}

	if resp.Errcode != 0 {
		return fmt.Errorf(resp.Errmsg)
	}

	return nil
}
