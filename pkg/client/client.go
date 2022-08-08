package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sovietscout/valorank/pkg/client/local"
	"github.com/sovietscout/valorank/pkg/content"
	"github.com/sovietscout/valorank/pkg/models"
	"github.com/sovietscout/valorank/pkg/riot"
)

type Client struct {
	IsRunning bool
	State     models.State

	Riot riot.IRiotClient

	stateCh chan models.State
	matchCh chan models.Match

	ctx    context.Context
	cancel context.CancelFunc
}

func NewClient(stateCh chan models.State, matchCh chan models.Match) *Client {
	return &Client{State: models.OFFLINE, stateCh: stateCh, matchCh: matchCh}
}

func (c *Client) Start() {
	log.Println("valorank: Client started")

	c.IsRunning = true
	c.ctx, c.cancel = context.WithCancel(context.Background())

	go func() {
		local.Client = local.NewClient(lockfileData())
		riot.SetVars(getRiotVars())

		content.SetContent()
		riot.SetCurrentSeason()

		c.ClientStateChangeLoop()
	}()
}

func (c *Client) Stop() {
	c.cancel()
	c.setState(models.OFFLINE)
	c.IsRunning = false

	log.Println("valorank: Client stopped")
}

// Checks every 1 second(s) to see if state has updated
func (c *Client) ClientStateChangeLoop() {
	log.Println("valorank: State change loop started")

	ticker := time.NewTicker(time.Second)

	L:
	for {
		select {
		case <- ticker.C:
			currPres := c.getPresence()

			if currPres == models.OFFLINE {
				continue
			}

			ticker.Stop()
			break L

		case <- c.ctx.Done():
			ticker.Stop()
			return
		}
	}

	conn, err := local.Client.InitWS()
	if err != nil {
		log.Println("ws init:", err)
		return
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.BinaryMessage, []byte(`[5, "OnJsonApiEvent_chat_v4_presences"]`))
	if err != nil {
		log.Println("ws write msg:", err)
	}

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("ws read:", err)
				break
			}

			if len(msg) == 0 {
				continue
			}

			var resp []json.RawMessage
			if err := json.Unmarshal(msg, &resp); err != nil {
				log.Println("json unmarshal:", err)
			}

			var respType string
			if err := json.Unmarshal(resp[1], &respType); err != nil {
				log.Println("json unmarshal:", err)
			}

			switch(respType) {
			case "OnJsonApiEvent_chat_v4_presences":
				var data WSChat4Presences
				if err := json.Unmarshal(resp[2], &data); err != nil {
					log.Println("json unmarshal:", err)
				}

				previousState := c.State

				for _, presence := range data.Data.Presences {
					if presence.Product == "valorant" && presence.Puuid == riot.UserPUUID {
						if data.EventType == "Delete" {
							c.setState(models.OFFLINE)

							log.Println("valorank: Game closed")
							return
						}

						newState := getStateFromPresence(&presence)

						if newState != models.MENU && newState == previousState {
							continue
						}

						c.setState(newState)
						previousState = newState

						break
					}
				}
			}
		}
	}()

	<- c.ctx.Done()

	err = conn.WriteMessage(
		websocket.CloseMessage, 
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	)

	if err != nil {
		log.Println("ws close:", err)
	}

	local.Client = nil
}

func (c *Client) getPresence() models.State {
	state := models.OFFLINE

	resp, err := local.Client.GET("/chat/v4/presences")
	if err != nil {
		log.Println("client get:", err)
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
			state = getStateFromPresence(&presence)
			break
		}
	}

	return state
}

func (c *Client) setState(state models.State) {
	log.Println("valorank: State set:", state)

	c.State = state
	c.stateCh <- state

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
	c.matchCh <- c.Riot.GetMatch()
	log.Println("valorank: Match received")
}

func getStateFromPresence(presence *riot.PresencesData) models.State {
	private_bytes, _ := base64.StdEncoding.DecodeString(presence.Private)

	data := new(riot.PresencesPrivate)
	json.Unmarshal(private_bytes, data)

	var state models.State

	switch data.SessionLoopState {
	case "MENUS":
		state = models.MENU
	case "PREGAME":
		state = models.PREGAME
	case "INGAME":
		state = models.INGAME
	}

	return state
}

type SessionResp struct {
	Puuid string `json:"puuid"`
}

type FetchSessionResp struct {
	LaunchConfiguration struct {
		Arguments []string `json:"arguments"`
	} `json:"launchConfiguration"`
}

type WSChat4Presences struct {
	Data struct {
		Presences []riot.PresencesData `json:"presences"`
	} `json:"data"`
	EventType string `json:"eventType"`
}