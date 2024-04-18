package creds

import (
	"encoding/json"
	"io/ioutil"
)

type TGCreds struct {
	TgAPIKey  string `json:"tgAPIkey"`
	TgChannel string `json:"tgChannel"`
}

type Credentials struct {
	Telegram TGCreds `json:"Telegram"`
}

var Loaded Credentials

func LoadCreds() (Credentials, error) {

	var c Credentials
	data, err := ioutil.ReadFile("./creds/creds.json")
	if err != nil {
		return Credentials{}, err
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		return Credentials{}, err
	}
	Loaded = c

	return c, nil
}
