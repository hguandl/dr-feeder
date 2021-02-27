package watcher

type indexData struct {
	Data struct {
		UserInfo struct {
			ScreenName string `json:"screen_name"`
		} `json:"userInfo"`
		TabsInfo struct {
			Tabs []struct {
				TabType     string `json:"tab_type"`
				Containerid string `json:"containerid"`
			} `json:"tabs"`
		} `json:"tabsInfo"`
	} `json:"data"`
}

type mblog struct {
	CreatedAt string                 `json:"created_at"`
	ID        string                 `json:"id"`
	Text      string                 `json:"text"`
	PicURL    string                 `json:"original_pic,omitempty"`
	PageInfo  map[string]interface{} `json:"page_info,omitempty"`
}

type cardData struct {
	Data struct {
		Cards []struct {
			CardType int   `json:"card_type"`
			Mblog    mblog `json:"mblog,omitempty"`
		} `json:"cards"`
	} `json:"data"`
}

type pageInfo struct {
	Type    string
	PagePic struct {
		URL string
	} `mapstructure:"page_pic"`
	PageURL string `mapstructure:"page_url"`
}
