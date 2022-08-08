package client

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sovietscout/valorank/pkg/client/local"
	"github.com/sovietscout/valorank/pkg/models"
)

var (
	doOnce sync.Once

	id string
	region models.Region
)

func lockfileData() (string, string, string) {
	data_bytes, _ := os.ReadFile(filepath.Join(os.Getenv("LOCALAPPDATA"), "Riot Games/Riot Client/Config/lockfile"))
	lockfile := strings.Split(string(data_bytes), ":")

	if len(lockfile) != 5 {
		time.Sleep(250 * time.Millisecond)
		return lockfileData()
	}

	return lockfile[2], lockfile[3], lockfile[4]
}

func getRiotVars() (string, models.Region) {
	doOnce.Do(func() {
		id = getID()
		region = getRegion()	
	})

	return id, region
}

func getID() string {
	resp, err := local.Client.GET("/chat/v1/session")
	if err != nil {
		time.Sleep(250 * time.Millisecond)
		return getID()		
	}

	data := new(SessionResp)
	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		time.Sleep(250 * time.Millisecond)
		return getID()
	}

	if data.Puuid == "" {
		time.Sleep(250 * time.Millisecond)
		return getID()
	}

	return data.Puuid
}

func getRegion() models.Region {
	resp, err := local.Client.GET("/product-session/v1/external-sessions")
	if err != nil {
		time.Sleep(250 * time.Millisecond)
		return getRegion()
	}

	data := map[string]FetchSessionResp{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		time.Sleep(250 * time.Millisecond)
		return getRegion()
	}

	var region models.Region

	for _, val := range data {
		for _, arg := range val.LaunchConfiguration.Arguments {
			if strings.HasPrefix(arg, "-ares-deployment") {
				region = models.GetRegion(arg[strings.Index(arg, "=") + 1:])
				break
			}
		}
	}

	return region
}