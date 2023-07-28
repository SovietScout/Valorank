package content

import (
	"encoding/json"
	"log"
	"net/http"
)

func SetContent() {
	if version, err := getClientVersion(); err == nil {
		ClientVersion = version
	} else {
		log.Println("There has been an error setting ClientVersion:", err)
	}
}

func getClientVersion() (string, error) {
	ver, err := http.Get("https://valorant-api.com/v1/version")
	if err != nil {
		return "", err
	}

	verData := new(VersionResp)
	json.NewDecoder(ver.Body).Decode(verData)

	return verData.Data.RiotClientVersion, nil
}

type VersionResp struct {
	Data struct {
		RiotClientVersion string `json:"riotClientVersion"`
	} `json:"data"`
}

type Data struct {
	Name   string
	Colour string
}
