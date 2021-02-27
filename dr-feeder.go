package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/hguandl/dr-feeder/v2/common"
	"github.com/hguandl/dr-feeder/v2/notifier"
	"github.com/hguandl/dr-feeder/v2/watcher"
)

// Version is current `git describe --tags` infomation.
var Version string = "v2.0.0"

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
	printVersion := flag.Bool("V", false, "Print current version")
	debugMode := flag.Bool("d", false, "Debug with fake server")
	pathPtr := flag.String("c", "config.yaml", "Configuration file")
	flag.Parse()

	if *printVersion {
		fmt.Printf("dr-feeder %s\n", Version)
		return
	}

	if *debugMode {
		println("Running on debug mode...")
	}

	config, err := LoadConfig(*pathPtr)
	if err != nil {
		log.Fatal(err)
	}

	notifiers, err := notifier.ParseNotifiers(config.Notifiers)
	if err != nil {
		log.Fatal(err)
	}

	watchers, err := watcher.ParseWatchers(config.Watchers, *debugMode)
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan common.NotifyPayload)

	for _, watcher := range watchers {
		go watch(watcher, ch)
	}

	go consume(ch, notifiers)

	select {}
}
