package watcher

type akWebAPIPayload struct {
	Code int           `json:"code"`
	Data akWebNewsList `json:"data"`
}

type akWebNewsList struct {
	List  []akWebNewsData `json:"list"`
	Total int             `json:"total"`
	End   bool            `json:"end"`
}

type akWebNewsData struct {
	Cid         string `json:"cid"`
	Tab         string `json:"tab"`
	Sticky      bool   `json:"sticky"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	DisplayTime int    `json:"displayTime"`
	Cover       string `json:"cover"`
	ExtraCover  string `json:"extraCover"`
	Brief       string `json:"brief"`
}
