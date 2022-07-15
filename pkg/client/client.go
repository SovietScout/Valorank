package client

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"strings"

	"github.com/sovietscout/valorank/pkg/content"
	"github.com/sovietscout/valorank/pkg/models"
	"github.com/sovietscout/valorank/pkg/riot"
)

type Client struct {
	State  State
	Region string
	PUUID  string

	ClientStateLoopOn bool
	PlayerChan chan ([]*models.Player)

	Riot	riot.RiotClient
}

func NewClient(playerChan chan []*models.Player) *Client {
	return &Client{State: OFFLINE, PlayerChan: playerChan}
}

func (c *Client) ClientReadyLoop(readyChan chan struct{}) {
	defer close(readyChan)

	ticker := time.NewTicker(time.Second)
	for range ticker.C {

		{
			// Set Riot Client
			data_bytes, err := os.ReadFile(filepath.Join(os.Getenv("LOCALAPPDATA"), "Riot Games/Riot Client/Config/lockfile"))
			if (err != nil) {
				continue
			}

			lockfile := strings.Split(string(data_bytes), ":")

			riot.Local = riot.NewNetCL(
				lockfile[2],
				lockfile[3],
				lockfile[4],
			)
		}

		{
			// Set PUUID
			if resp, err := riot.Local.GET("/chat/v1/session"); err == nil {
				data := new(SessionResp)
				json.NewDecoder(resp.Body).Decode(data)
	
				c.PUUID = data.Puuid
			}
		}

		{
			// Set region
			if resp, err := riot.Local.GET("/product-session/v1/external-sessions"); err == nil {
				data := map[string]FetchSessionResp{}
				json.NewDecoder(resp.Body).Decode(&data)

				for _, val := range data {
					for _, arg := range val.LaunchConfiguration.Arguments {
						if strings.HasPrefix(arg, "-ares-deployment") {
							c.Region = arg[strings.Index(arg, "=") + 1:]
						}
					}
				}
			}
		}

		if riot.Local != nil && c.PUUID != "" && c.Region != "" {
			readyChan <- struct{}{}

			content.SetContent()
			riot.SetVars(c.PUUID, c.Region)
			riot.SetCurrentSeason()

			ticker.Stop()
		}
	}
}

// Checks every 1 second(s) to see if state has updated
func (c *Client) ClientStateChangeLoop(ret chan State) {
	ticker := time.NewTicker(time.Second)
	previousState := c.State

	c.ClientStateLoopOn = true

	for range ticker.C {
		var currentState = c.getPresence()

		if currentState != previousState {
			c.setState(currentState)
			previousState = currentState

			ret <- currentState
		}
	}
}

func (c *Client) getPresence() State {
	if riot.Local == nil {
		return OFFLINE
	}

	resp, err := riot.Local.GET("/chat/v4/presences")
	if err != nil {
		return OFFLINE
	}

	data := new(riot.PresencesResp)
	json.NewDecoder(resp.Body).Decode(data)
	resp.Body.Close()

	if len(data.Presences) == 0 {
		return OFFLINE
	}

	state := OFFLINE

	for _, presence := range data.Presences {
		if presence.Product == "valorant" && presence.Puuid == c.PUUID {
			private_bytes, _ := base64.StdEncoding.DecodeString(presence.Private)

			data := new(riot.PresencesPrivate)
			json.Unmarshal(private_bytes, data)

			switch data.SessionLoopState {
			case "MENUS":
				state = MENU
			case "PREGAME":
				state = PREGAME
			case "INGAME":
				state = INGAME
			}

			break
		}
	}

	return state
}

func (c *Client) setState(state State) {
	c.State = state

	switch state {
	case MENU:
		c.Riot = new(riot.Menu)
	case PREGAME:
		c.Riot = new(riot.Pregame)
	case INGAME:
		c.Riot = new(riot.Ingame)
	default:
		c.Riot = nil	// When Offline
		return
	}

	go c.Riot.GetPlayers(c.PlayerChan)
}

type SessionResp struct {
	Puuid     string `json:"puuid"`
}

type FetchSessionResp struct {
	LaunchConfiguration struct {
		Arguments []string `json:"arguments"`
	} `json:"launchConfiguration"`
}