package local

import (
	"encoding/json"
	"net/http"

	"github.com/sovietscout/valorank/pkg/content"
	"golang.org/x/sync/singleflight"
)

var requestGroup singleflight.Group

func (n *LocalClient) GetRiotHeaders() http.Header {
	v, err, _ := requestGroup.Do("local", func() (interface{}, error) {
		return n.generateRiotHeaders()
	})

	if err != nil {
		return nil
	}

	return v.(http.Header)
}

func (n *LocalClient) generateRiotHeaders() (http.Header, error) {
	resp, err := n.GET("/entitlements/v1/token")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data := new(EntitlementResp)
	json.NewDecoder(resp.Body).Decode(data)

	return http.Header{
		"Authorization":           {"Bearer " + data.AccessToken},
		"X-Riot-Entitlements-JWT": {data.Token},
		"X-Riot-ClientPlatform":   {content.ClientPlatform},
		"X-Riot-ClientVersion":    {content.ClientVersion},
		"User-Agent":              {"ShooterGame/13 Windows/10.0.19043.1.256.64bit"},
	}, nil
}

type EntitlementResp struct {
	AccessToken string `json:"accessToken"`
	Token       string `json:"token"`
}