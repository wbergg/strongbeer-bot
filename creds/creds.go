package creds

import (
	"encoding/json"
	"io/ioutil"
)

type TGCreds struct {
	tgAPIKey  string
	tgChannel string
}

type Credentials struct {
	Telegram TGCreds
}

var Loaded Credentials

func LoadCreds() (Credentials, error) {

	var c Credentials
	data, err := ioutil.ReadFile("./creds/creds.json")
	if err != nil {
		return Credentials{}, err
	}

	json.Unmarshal(data, &c)

	Loaded = c

	return c, nil
}
