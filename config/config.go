package config

import (
	"encoding/json"
	"os"
)

type TGCreds struct {
	TgAPIKey  string `json:"tgAPIkey"`
	TgChannel string `json:"tgChannel"`
}

type Config struct {
	Telegram    TGCreds           `json:"Telegram"`
	Spreadsheet string            `json:"SpreadsheetID"`
	Sheetname   string            `json:"Sheetname"`
	Users       map[string]string `json:"UserMap"`
}

var Loaded Config

func LoadConfig(configFile string) (Config, error) {

	var c Config
	data, err := os.ReadFile(configFile)
	if err != nil {
		return Config{}, err
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		return Config{}, err
	}
	Loaded = c

	return c, nil
}
