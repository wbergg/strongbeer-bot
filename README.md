# strongbeer-bot
![strongbeer-logo](https://assets.untappd.com/site/beer_logos_hd/beer-55370_24b70_hd.jpeg)

## About The Project

A telegram bot to keep track of whether it's Starkölsmåndag or not. Also featuring various other functions collected from a Google Sheet.

More information about Starkölsmåndag can be found at [this site.](https://starkölsmåndag.se)

### Prerequisites
To be able to use the Google API you need to create an service account on your Google Console and download the json-key and put in the config/ dir.

## Getting Started

Config and Telegram information are defined in config/config.json file using the following template:

Telegram APIkey, Channel, SpreadsheetID and a UserMap:

```
{
    "Telegram": {
		"tgAPIkey": "xxx",
		"tgChannel": "xxx"
	}
      "SpreadsheetID": "xxx",
      "UserMap": {
		"user1": "id",
            "user2": "id"
      }

}
```

## Running

```
go run main.go
```

### DEBUG mode
```
  -config-file string
        Absolute path for config-file (default "./config/config.json")
  -debug
        true/false - Turns on debug for telegram
  -google-creds-file string
        Absolute path for Google creds-file (default "./config/google.json")
  -stdout
        true/false - Turns on stdout rather than sending to telegra
  -telegram-test
        Sends a test message to specified telegram channel
```