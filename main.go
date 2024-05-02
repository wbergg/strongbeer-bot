package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/wbergg/strongbeer-bot/creds"
	"github.com/wbergg/strongbeer-bot/tele"
)

var announce bool = false
var debug_telegram *bool
var debug_stdout *bool
var nextMondayDate time.Time

func main() {
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

	channel, err := strconv.ParseInt(credentials.Telegram.TgChannel, 10, 64)
	if err != nil {
		log.Error(err)
		panic("Could not convert Telegram channel to int64")
	}
	tg := tele.New(credentials.Telegram.TgAPIKey, channel, *debug_telegram, *debug_stdout)
	tg.Init(*debug_telegram)

	// Set date for next monday
	nextMonday()

	// Initiate read message function
	go readMessage(tg)

	// Loop forever
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
			go mondayReminder(tele, 12)
			go mondayReminder(tele, 18)
			go mondayReminder(tele, 21)
			nextMondayDate = t.AddDate(0, 0, 7)
		} else {
			return
		}
	}
	if t.Weekday() != time.Monday {
		announce = false
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
		log.Error(err)
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
				"/status - Status of bot\n"
		case "status":
			message = "I'm ok."
		case "som":
			t := time.Now()
			timeUntilMsg := "\n\nTime until next STARKÖLSMÅNDAG is: " + timeUntilFormatted(t, nextMondayDate)
			if t.Weekday() == time.Monday {
				message = "\xF0\x9F\x8D\xBB YES! IT'S STARKÖLSMÅNDAG! \xF0\x9F\x8D\xBB"
			} else if t.Weekday() == time.Tuesday {
				message = "No, it's not.\n\n" +
					"However, it's \xF0\x9F\x8D\xB8 GIN & TONIC tisdag! \xF0\x9F\x8D\xB8" + timeUntilMsg
			} else if t.Weekday() == time.Wednesday {
				message = "No, it's not.\n\n" +
					"However, it's \xF0\x9F\x8E\x89 BERGFEST TAG! \xF0\x9F\x8E\x89" + timeUntilMsg
			} else if t.Weekday() == time.Thursday {
				message = "No, it's not.\n\n" +
					"However, it's \xF0\x9F\xA5\x83 WHISKEY TORSDAG! \xF0\x9F\xA5\x83" + timeUntilMsg
			} else if t.Weekday() == time.Saturday {
				message = "No, it's not.\n\n" +
					"However, it's \xF0\x9F\x8D\xB7 VIN LÖRDAG! \xF0\x9F\x8D\xB7" + timeUntilMsg
			} else {
				message = "No, it's not." + timeUntilMsg
			}

		default:
			message = "I don't know that command"
		}

		if _, err := tele.SendM(message); err != nil {
			log.Panic(err)
		}
	}
}

func timeUntilFormatted(a time.Time, b time.Time) string {
	d := b.Sub(a).Round(time.Second) // Round to nearest second

	if d < 0 {
		d = -d
	}

	const day = 24 * time.Hour

	if d < day {
		return d.String()
	}

	n := int(d / day)
	d -= time.Duration(n) * day

	res := fmt.Sprintf("%dd%s", n, d.Round(time.Second))

	return res
}

func nextMonday() {
	t := time.Now()
	daysUntilMonday := (7 - int(t.Weekday()) + 1) % 7
	nextMondayDate = t.AddDate(0, 0, daysUntilMonday)
	nextMondayDate = time.Date(nextMondayDate.Year(), nextMondayDate.Month(), nextMondayDate.Day(), 0, 0, 0, 0, time.Local)
	// Debug
	//fmt.Println(nextMondayDate)
}
