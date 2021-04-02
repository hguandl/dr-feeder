package watcher

type sirenAPIPayload struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

type sirenListData struct {
	List []sirenNewsData `mapstructure:"list"`
	End  bool            `mapstructure:"end"`
}

type sirenNewsData struct {
	Cid     string `mapstructure:"cid"`
	Title   string `mapstructure:"title"`
	Cate    int    `mapstructure:"cate"`
	Author  string `mapstructure:"author,omitempty"`
	Content string `mapstructure:"content,omitempty"`
	Date    string `mapstructure:"date"`
}
