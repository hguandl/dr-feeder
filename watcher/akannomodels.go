package watcher

type announce struct {
	AnnounceID string `json:"announceId"`
	Title      string `json:"title"`
	IsWebURL   bool   `json:"isWebUrl"`
	WebURL     string `json:"webUrl"`
	Day        int    `json:"day"`
	Month      int    `json:"month"`
	Group      string `json:"group"`
}

type announceMeta struct {
	FocusAnnounceID string     `json:"focusAnnounceId"`
	AnnounceList    []announce `json:"announceList"`
}
