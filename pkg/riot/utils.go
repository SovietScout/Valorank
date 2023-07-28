package riot

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sovietscout/valorank/pkg/local"
	"github.com/sovietscout/valorank/pkg/content"
	"github.com/sovietscout/valorank/pkg/models"
)

func SetVars(userPUUID string, region models.Region) {
	UserPUUID = userPUUID
	Region = region

	log.Println("Riot vars set")
}

func SetCurrentSeason() error {
	if content.CurrentSeasonID != "" {
		return nil
	}

	req, _ := http.NewRequest(http.MethodGet, GetSharedURL("/content-service/v3/content"), nil)
	req.Header = local.Client.GetRiotHeaders()

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data := new(FetchContentResp)
	json.NewDecoder(resp.Body).Decode(data)

	for _, season := range data.Seasons {
		if season.IsActive && season.Type == "act" {
			content.CurrentSeasonID = season.ID
			break
		}
	}

	return nil
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
	return "https://glz-" + Region.Region + "-1." + Region.Shard + ".a.pvp.net" + endpoint
}

func GetPDURL(endpoint string) string {
	return "https://pd." + Region.Shard + ".a.pvp.net" + endpoint
}

func GetSharedURL(endpoint string) string {
	return "https://shared." + Region.Shard + ".a.pvp.net" + endpoint
}
