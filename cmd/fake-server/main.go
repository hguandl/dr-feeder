package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

var weiboTests = [6]string{
	"tests/weibo/01-mblog-with-article.json",
	"tests/weibo/02-mblog-with-video.json",
	"tests/weibo/03-mblog-with-text.json",
	"tests/weibo/04-mblog-with-tag-and-pic.json",
	"tests/weibo/05-retweeted.json",
	"tests/weibo/06-lottery.json",
}
var weiboIdx = 0

var akAnnoTests = [6]string{
	"tests/akanno/00-init.json",
	"tests/akanno/01-new-gacha.json",
	"tests/akanno/02-activity-end.json",
	"tests/akanno/03-placehold.json",
	"tests/akanno/04-dev-news.json",
	"tests/akanno/05-null.json",
}
var akAnnoIdx = 0

var sirenTests = [6]string{
	"tests/siren/00.json",
	"tests/siren/01.json",
	"tests/siren/750450.json",
}
var sirenIdx = 0

func weiboHandler(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadFile(weiboTests[weiboIdx])
	log.Printf("Deliverd %v\n", weiboTests[weiboIdx])
	weiboIdx = (weiboIdx + 1) % 6
	w.Write(data)
}

func akAnnoHandler(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadFile(akAnnoTests[akAnnoIdx])
	log.Printf("Deliverd %v\n", akAnnoTests[akAnnoIdx])
	akAnnoIdx = (akAnnoIdx + 1) % 6
	w.Write(data)
}

func sirenHandler(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadFile(sirenTests[sirenIdx])
	log.Printf("Deliverd %v\n", sirenTests[sirenIdx])
	sirenIdx = (sirenIdx + 1) % 3
	w.Write(data)
}

func main() {
	listenAddr := ":8088"
	http.HandleFunc("/weibo", weiboHandler)
	http.HandleFunc("/akanno", akAnnoHandler)
	http.HandleFunc("/siren", sirenHandler)

	log.Printf("Listen at %v\n", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
