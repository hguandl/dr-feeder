package wxmsgapp

import (
	"fmt"
	"time"
)

// WxAPIClient communicates with Work Wechat API server.
type WxAPIClient struct {
	AgentID     string    `json:"agentid" binding:"required"`
	ToUser      string    `json:"touser" binding:"required"`
	CorpID      string    `json:"corpid" binding:"required"`
	CorpSecret  string    `json:"corpsecret" binding:"required"`
	AccessToken string    `json:"-"`
	TokenUntil  time.Time `json:"-"`
}

type wxAPIResp struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func (client *WxAPIClient) String() string {
	return fmt.Sprintf("{%v %v %v %v %v %v}",
		client.AgentID,
		client.ToUser,
		client.CorpID,
		client.CorpSecret,
		client.AccessToken,
		client.TokenUntil,
	)
}
