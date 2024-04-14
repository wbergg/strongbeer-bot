package main

import (
	"flag"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/wbergg/strongbeer-bot/creds"

	"github.com/wbergg/bordershop-bot/tele"
)

func main() {
	// Enable bool debug flag
	debug := flag.Bool("debug", false, "Turns on debug mode and prints to stdout")

	// Load credentials
	credentials, err := creds.LoadSecrets()
	if err != nil {
		log.Error(err)
		panic("Could not load credentials, check creds/creds.json")
	}

	tg := tele.New(tgAPIKey, tgChannel, false, *debug)
	tg.Init(false)

	// Run before starting the timer
	mondayTimer(tg)
}

func mondayTimer(t *tele.Tele) {
	fmt.Println("test")
}
