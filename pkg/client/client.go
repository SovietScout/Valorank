package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"time"

	"github.com/sovietscout/valorank/pkg/content"
	"github.com/sovietscout/valorank/pkg/models"
	"github.com/sovietscout/valorank/pkg/riot"
)

type Client struct {
	IsRunning bool
	State     models.State

	Riot riot.RiotClient

	matchChan chan models.Match

	ctx    context.Context
	cancel context.CancelFunc
}

func NewClient(matchChan chan models.Match) *Client {
	return &Client{State: models.OFFLINE, matchChan: matchChan}
}

func (c *Client) Start(ret chan models.State) {
	log.Println("Client started")

	c.IsRunning = true
	c.ctx, c.cancel = context.WithCancel(context.Background())

	go func() {
		riot.Local = riot.NewNetCL(lockfileData())
		riot.SetVars(getRiotVars())

		content.SetContent()
		riot.SetCurrentSeason()

		c.ClientStateChangeLoop(ret)
	}()
}

func (c *Client) Stop() {
	c.cancel()
	c.setState(models.OFFLINE)
	c.IsRunning = false

	log.Println("Client stopped")
}

// Checks every 1 second(s) to see if state has updated
func (c *Client) ClientStateChangeLoop(ret chan models.State) {
	log.Println("State change loop started")

	ticker := time.NewTicker(time.Second)
	previousState := c.State

	for {
		select {
		case <-ticker.C:
			if currentState := c.getPresence(); currentState != previousState {
				c.setState(currentState)
				previousState = currentState

				ret <- currentState
			}

		case <-c.ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (c *Client) getPresence() models.State {
	state := models.OFFLINE

	resp, err := riot.Local.GET("/chat/v4/presences")
	if err != nil {
		log.Println("Err in Get:", err)
		return state
	}
	defer resp.Body.Close()

	data := new(riot.PresencesResp)
	json.NewDecoder(resp.Body).Decode(data)

	if len(data.Presences) == 0 {
		return state
	}

	for _, presence := range data.Presences {
		if presence.Product == "valorant" && presence.Puuid == riot.UserPUUID {
			private_bytes, err := base64.StdEncoding.DecodeString(presence.Private)
			if err != nil {
				log.Println("Err in decode string")
			}

			data := new(riot.PresencesPrivate)
			json.Unmarshal(private_bytes, data)

			switch data.SessionLoopState {
			case "MENUS":
				state = models.MENU
			case "PREGAME":
				state = models.PREGAME
			case "INGAME":
				state = models.INGAME
			default:
				log.Println(presence)	// log whole presence if offline
			}

			break
		}
	}

	return state
}

func (c *Client) setState(state models.State) {
	log.Printf("State set: %s\n", state)
	c.State = state

	switch state {
	case models.MENU:
		c.Riot = new(riot.Menu)
	case models.PREGAME:
		c.Riot = new(riot.Pregame)
	case models.INGAME:
		c.Riot = new(riot.Ingame)
	default:
		c.Riot = nil // When Offline
		return
	}

	go c.GetMatch()
}

// Essentially a wrapper
func (c *Client) GetMatch() {
	log.Println("Getting match")
	c.matchChan <- c.Riot.GetMatch()
	log.Println("Match received")
}

type SessionResp struct {
	Puuid string `json:"puuid"`
}

type FetchSessionResp struct {
	LaunchConfiguration struct {
		Arguments []string `json:"arguments"`
	} `json:"launchConfiguration"`
}
