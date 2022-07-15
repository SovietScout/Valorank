package riot

import (
	"encoding/json"
	"net/http"

	"github.com/sovietscout/valorank/pkg/content"
	"golang.org/x/sync/singleflight"
)

var requestGroup singleflight.Group

func (n *NetCL) GetRiotHeaders() http.Header {
	v, err, _ := requestGroup.Do("local", func() (interface{}, error) {
		return n.generateRiotHeaders()
	})

	if err != nil {
		return nil
	}

	return v.(http.Header)
}

func (n *NetCL) generateRiotHeaders() (http.Header, error) {
	resp, err := n.GET("/entitlements/v1/token")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data := new(EntitlementResp)
	json.NewDecoder(resp.Body).Decode(data)

	return http.Header{
		"Authorization":			{"Bearer " + data.AccessToken},
		"X-Riot-Entitlements-JWT":	{data.Token},
		"X-Riot-ClientPlatform":	{content.ClientPlatform},
		"X-Riot-ClientVersion": 	{content.ClientVersion},
		"User-Agent":           	{"ShooterGame/13 Windows/10.0.19043.1.256.64bit"},
	}, nil
}

func SetVars(userPUUID, region string) {
	UserPUUID = userPUUID
	Region = region

	SetCurrentSeason()
}

func SetCurrentSeason() {
	if content.CurrentSeasonID != "" {
		return
	}

	req, _ := http.NewRequest(http.MethodGet, GetSharedURL("/content-service/v3/content"), nil)
	req.Header = Local.GetRiotHeaders()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return // Don't know what the hell to do
	}
	defer resp.Body.Close()

	data := new(FetchContentResp)
	json.NewDecoder(resp.Body).Decode(data)

	for _, season := range data.Seasons {
		if season.IsActive && season.Type == "act" {
			content.CurrentSeasonID = season.ID
		}
	}
}

type EntitlementResp struct {
	AccessToken string `json:"accessToken"`
	Token       string `json:"token"`
}

type FetchContentResp struct {
	Seasons []struct {
		ID       string `json:"ID"`
		Type     string `json:"Type"`
		IsActive bool   `json:"IsActive"`
	} `json:"Seasons"`
}

func GetGLZURL(endpoint string) string {
	return "https://glz-" + Region + "-1." + Region + ".a.pvp.net" + endpoint
}

func GetPDURL(endpoint string) string {
	return "https://pd." + Region + ".a.pvp.net" + endpoint
}

func GetSharedURL(endpoint string) string {
	return "https://shared." + Region + ".a.pvp.net" + endpoint
}


