package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/hguandl/dr-feeder/v2/common"
	"github.com/hguandl/dr-feeder/v2/notifier"
	"github.com/hguandl/dr-feeder/v2/watcher"
)

func consume(ch chan common.NotifyPayload, notifiers []notifier.Notifier) {
	for {
		pl := <-ch
		log.Printf("Received \"%s\":\n==========\n%s\n==========",
			pl.Title, pl.Body)

		for _, ntf := range notifiers {
			go ntf.Push(pl)
		}
	}
}

func watch(watcher watcher.Watcher, ch chan common.NotifyPayload) {
	for {
		waitSec := rand.Intn(6) + 6
		watcher.Produce(ch)
		time.Sleep(time.Duration(waitSec) * time.Second)
	}
}

func main() {
	pathPtr := flag.String("c", "config.yaml", "Configuration file")
	flag.Parse()

	notifiers, err := ParseConfig(*pathPtr)
	if err != nil {
		log.Fatal(err)
	}

	weibo, err := watcher.NewWeiboWatcher(6279793937)
	if err != nil {
		log.Fatal(err)
	}

	anAnno, err := watcher.NewAkAnnounceWatcher()
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan common.NotifyPayload)

	go watch(weibo, ch)
	go watch(anAnno, ch)

	go consume(ch, notifiers)

	select {}
}
