package client

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sovietscout/valorank/pkg/riot"
)

var (
	doOnce sync.Once

	id, region string
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

func getRiotVars() (string, string) {
	doOnce.Do(func() {
		id = getID()
		region = getRegion()	
	})

	return id, region
}

func getID() string {
	resp, err := riot.Local.GET("/chat/v1/session")
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

func getRegion() string {
	resp, err := riot.Local.GET("/product-session/v1/external-sessions")
	if err != nil {
		time.Sleep(250 * time.Millisecond)
		return getRegion()
	}

	data := map[string]FetchSessionResp{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		time.Sleep(250 * time.Millisecond)
		return getRegion()
	}

	var region string

	for _, val := range data {
		for _, arg := range val.LaunchConfiguration.Arguments {
			if strings.HasPrefix(arg, "-ares-deployment") {
				region = arg[strings.Index(arg, "=") + 1:]
				break
			}
		}
	}

	return region
}