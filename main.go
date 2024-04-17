package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/wbergg/strongbeer-bot/creds"
	"github.com/wbergg/strongbeer-bot/tele"
)

var announce bool = false
var debug_telegram *bool
var debug_stdout *bool
var nextMonday time.Time

func main() {
	// temp shit
	apikey := os.Getenv("tgAPIKey")
	channel, _ := strconv.ParseInt(os.Getenv("tgChannel"), 10, 64)

	fmt.Println(channel, apikey)

	// Enable bool debug flag
	debug_telegram = flag.Bool("debug", false, "true/false - Turns on debug for telegram")
	debug_stdout = flag.Bool("stdout", false, "true/false - Turns on stdout rather than sending to telegra")
	flag.Parse()

	// Load credentials
	credentials, err := creds.LoadCreds()
	if err != nil {
		log.Error(err)
		panic("Could not load credentials, check creds/creds.json")
	}
	fmt.Println(credentials.Telegram)
	//tg := tele.New(credentials.Telegram.tgChannel, credentials.Telegram.tgAPIKey, false, *debug)
	tg := tele.New(apikey, channel, *debug_telegram, *debug_stdout)
	tg.Init(*debug_telegram)

	// Initiate read message function
	go readMessage(tg)

	for {
		mondayTimer(tg)
		time.Sleep(5 * time.Second)
	}

}

func mondayTimer(tele *tele.Tele) {

	t := time.Now()

	if t.Weekday() == time.Monday {
		if announce == false {
			tele.SendM("\xF0\x9F\x8D\xBA IT'S STARKÖLSMÅNDAG! \xF0\x9F\x8D\xBA")
			announce = true
			go mondayReminder(tele, 5)
			go mondayReminder(tele, 10)
			nextMonday = t.AddDate(0, 0, 7)
			//fmt.Println("This the next monday is:", renewOn.Format(time.RFC822))
		} else {
			return
		}
		if t.Weekday() != time.Monday {
			announce = false
		}
	}
}

func mondayReminder(tele *tele.Tele, t time.Duration) {
	timer := time.NewTimer(t * time.Hour)
	<-timer.C
	tele.SendM("\xF0\x9F\x8D\xBA REMINDER - IT'S STARKÖLSMÅNDAG! \xF0\x9F\x8D\xBA")
}

func readMessage(tele *tele.Tele) {

	updates, err := tele.ReadM()
	if err != nil {
		panic(err)
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}
		var message string

		// Do something with the update
		switch update.Message.Command() {
		case "help":
			message = "Hi\n\n" +
				"I'm STRONGBEER-BOT, these are the current commands:\n\n" +
				"/som - Is it starkölsmåndag or not?\n" +
				"/status - Status\n"
		case "status":
			message = "I'm ok."
		case "som":
			t := time.Now()
			if t.Weekday() == time.Monday {
				message = "\xF0\x9F\x8D\xBB YES! IT'S STARKÖLSMÅNDAG! \xF0\x9F\x8D\xBB"
			} else if t.Weekday() == time.Tuesday {
				message = "No, it's not.\n\n" +
					"However, it's \xF0\x9F\x8D\xB8 GIN & TONIC tisdag! \xF0\x9F\x8D\xB8"
			} else if t.Weekday() == time.Wednesday {
				message = "No, it's not.\n\n" +
					"However, it's \xF0\x9F\x8E\x89 BERGFEST TAG! \xF0\x9F\x8E\x89"
			} else if t.Weekday() == time.Thursday {
				message = "No, it's not.\n\n" +
					"However, it's \xF0\x9F\xA5\x83 WHISKEY TORSDAG! \xF0\x9F\xA5\x83"
			} else if t.Weekday() == time.Saturday {
				message = "No, it's not.\n\n" +
					"However, it's \xF0\x9F\x8D\xB7 VIN LÖRDAG! \xF0\x9F\x8D\xB7"
			} else {
				message = "No, it's not.\n\n" +
					"The next Monday is: " +
					nextMonday.Format(time.RFC822)
			}

		default:
			message = "I don't know that command"
		}

		if _, err := tele.SendM(message); err != nil {
			log.Panic(err)
		}
	}
}
