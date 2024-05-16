# strongbeer-bot
![strongbeer-logo](https://assets.untappd.com/site/beer_logos_hd/beer-55370_24b70_hd.jpeg)

A telegram bot to keep track of whether it's Starkölsmåndag or not.

Credentials are defined in creds/creds.json file using the following template:

Telegram APIkey and channel definition:

```
{
    "Telegram": {
		"tgAPIkey": "xxx",
		"tgChannel": "xxx"
	}
}
```

## Running

```
go run main.go
```

### DEBUG mode
```
  -debug
        true/false - Turns on debug for telegram
  -stdout
        true/false - Turns on stdout rather than sending to telegram
```