package main

import (
	"flag"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/wbergg/strongbeer-bot/creds"
)

var announce bool = false

func main() {
	// Enable bool debug flag
	debug := flag.Bool("debug", false, "true/false - Turns on debug mode and prints to stdout")
	flag.Parse()

	// Load credentials
	credentials, err := creds.LoadCreds()
	if err != nil {
		log.Error(err)
		panic("Could not load credentials, check creds/creds.json")
	}
	fmt.Println(credentials.Telegram)
	//tg := tele.New(credentials.Telegram.tgChannel, credentials.Telegram.tgAPIKey, false, *debug)
	//tg.Init(false)

	// Run before starting the timer
	//mondayTimer(tg)
	for {
		if *debug {
			fmt.Println("true")
		} else {
			mondayTimer()
		}
		time.Sleep(5 * time.Second)
	}

}

// func mondayTimer(t *tele.Tele) {
func mondayTimer() {

	t := time.Now()

	if t.Weekday() == time.Monday {
		if announce == false {
			fmt.Println("Monday")
			announce = true
			go mondayReminer(12)
			go mondayReminer(18)
		} else {
			return
		}
		if t.Weekday() != time.Monday {
			announce = false
		}
	}

	fmt.Println("test")
}

func mondayReminer(d time.Duration) {
	timer := time.NewTimer(d * time.Second)

	<-timer.C
	fmt.Println("Timer fired!")

}

func readMessage() {

}
