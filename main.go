package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"github.com/wbergg/strongbeer-bot/config"
	"github.com/wbergg/strongbeer-bot/sheetservice"
	"github.com/wbergg/strongbeer-bot/tele"
)

var announce bool = false
var nextMondayDate time.Time
var debugStdout *bool

func main() {
	// Enable bool debug flag
	debugTelegram := flag.Bool("debug", false, "true/false - Turns on debug for telegram")
	debugStdout = flag.Bool("stdout", false, "true/false - Turns on stdout rather than sending to telegra")
	telegramTest := flag.Bool("telegram-test", false, "Sends a test message to specified telegram channel")
	configFile := flag.String("config-file", "./config/config.json", "Absolute path for config-file")
	googleCreds := flag.String("google-creds-file", "./config/google.json", "Absolute path for Google creds-file")
	flag.Parse()

	// Load credentials
	config, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Error(err)
		panic("Could not load credentials, check config file.")
	}

	channel, err := strconv.ParseInt(config.Telegram.TgChannel, 10, 64)
	if err != nil {
		log.Error(err)
		panic("Could not convert Telegram channel to int64")
	}
	tg := tele.New(config.Telegram.TgAPIKey, channel, *debugTelegram, *debugStdout)
	tg.Init(*debugTelegram)

	//Set up data logging
	f, err := os.OpenFile("strongbeer-log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	logrus.SetOutput(f)

	// Init Google sheets
	s, err := sheetservice.InitSheetService(config.Spreadsheet, *googleCreds)
	if err != nil {
		log.Fatal(err)
	}

	if *telegramTest {
		tg.SendM("DEBUG: strongbeer-bot test message")
		// End program after sending message
		os.Exit(0)
	} else {

		// Set date for next monday
		nextMonday()

		// Initiate read message function
		go readMessage(tg, s)

		// Loop forever
		for {
			mondayTimer(tg)
			time.Sleep(5 * time.Second)
		}
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
			if *debugStdout {
				fmt.Println("Next Monday date is: ", nextMondayDate.String())
			}
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

func readMessage(tele *tele.Tele, s *sheetservice.SheetService) {

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

		// Switch/case for commands
		switch update.Message.Command() {
		case "help":
			message = "Hi\n\n" +
				"I'm STRONGBEER-BOT, these are the current commands:\n\n" +
				"/som - Is it starkölsmåndag or not?\n" +
				"/botstatus - Status of bot\n" +
				"/scoreboard - Shows leader board\n" +
				"/top3 - Shows top3 SÖMers\n" +
				"/checkin - Show your checkin for current week\n"

		case "botstatus":
			message = "I'm ok."

		case "status":
			data, err := s.GetSheetUserData(update.Message.From.UserName)
			if err != nil {
				log.Errorf("Error getting sheet data: %v", err)
				message = "An error occurred while retrieving your status data."
			} else {
				message = data
			}

		case "checkin":
			data, err := s.GetSheetUserCheckin(update.Message.From.UserName)
			if err != nil {
				log.Errorf("Error getting sheet data: %v", err)
				message = "An error occurred while retrieving your checkin data."
			} else {
				message = data
			}

		case "scoreboard":
			data, err := s.GetSheetTopList(false)
			if err != nil {
				log.Errorf("Error getting sheet data: %v", err)
				message = "An error occurred while retrieving sheet data."
			} else {
				message = data
			}

		case "top3":
			data, err := s.GetSheetTopList(true)
			if err != nil {
				log.Errorf("Error getting sheet data: %v", err)
				message = "An error occurred while retrieving sheet data."
			} else {
				message = data
			}

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
	if *debugStdout {
		fmt.Println("Next Monday date is: ", nextMondayDate.String())
	}
}
