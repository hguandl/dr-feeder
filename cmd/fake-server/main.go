package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

var weiboTests = [4]string{
	"tests/weibo/01-mblog-with-article.json",
	"tests/weibo/02-mblog-with-video.json",
	"tests/weibo/03-mblog-with-text.json",
	"tests/weibo/04-mblog-with-tag-and-pic.json",
}

var weiboIdx = 0

func handler(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadFile(weiboTests[weiboIdx])
	log.Printf("Deliverd %v\n", weiboTests[weiboIdx])
	weiboIdx = (weiboIdx + 1) % 4
	w.Write(data)
}

func main() {
	listenAddr := ":8088"
	http.HandleFunc("/weibo", handler)

	log.Printf("Listen at %v\n", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
