package content

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

func SetContent() {
	ClientVersion = getClientVersion()
}

func getClientVersion() string {
	ver, err := http.Get("https://valorant-api.com/v1/version")
	if err != nil {
		return ""
	}

	verData := new(VersionResp)
	json.NewDecoder(ver.Body).Decode(verData)

	return verData.Data.RiotClientVersion
}

type VersionResp struct {
	Data struct {
		RiotClientVersion string `json:"riotClientVersion"`
	} `json:"data"`
}

type AgentData struct {
	Name   string
	Colour string
}

type RankData struct {
	Name string
	Colour string
}
