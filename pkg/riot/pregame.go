package riot

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"

	"github.com/sovietscout/valorank/pkg/local"
	"github.com/sovietscout/valorank/pkg/models"
)

var ErrPregameIndexOutOfRange = errors.New("pregame: Team slice empty, likely moved to Game")

type Pregame struct {
	sync.Mutex
}

func (r *Pregame) GetMatch() models.Match {
	r.Lock()
	defer r.Unlock()

	match := models.Match{State: models.PREGAME}

	// Current Pregame ID
	reqID, _ := http.NewRequest(http.MethodGet, GetGLZURL("/pregame/v1/players/"+UserPUUID), nil)
	reqID.Header = local.Client.GetRiotHeaders()

	respID, err := client.Do(reqID)
	if err != nil {
		match.Err = err
		return match
	}
	defer respID.Body.Close()

	dataID := new(GameFetchPlayerResp)
	json.NewDecoder(respID.Body).Decode(dataID)

	// Pregame details
	req, _ := http.NewRequest(http.MethodGet, GetGLZURL("/pregame/v1/matches/"+dataID.MatchID), nil)
	req.Header = local.Client.GetRiotHeaders()

	resp, err := client.Do(req)
	if err != nil {
		match.Err = err
		return match
	}
	defer resp.Body.Close()

	data := new(PregameFetchMatchResp)
	json.NewDecoder(resp.Body).Decode(data)

	match.GamePodID = data.GamePodID

	if len(data.Teams) == 0 {
		match.Err = ErrPregameIndexOutOfRange
		return match
	}

	playerRecv := make(chan *models.Player, len(data.Teams[0].Players))
	var wg sync.WaitGroup

	for _, player := range data.Teams[0].Players {
		wg.Add(1)

		go func(player PlayerResp) {
			defer wg.Done()

			playerLevel := -1
			if !player.PlayerIdentity.HideAccountLevel {
				playerLevel = player.PlayerIdentity.AccountLevel
			}

			p := &models.Player{
				SubjectID: player.Subject,
				Level:     playerLevel,
				Ally:      true,
				Incognito: player.PlayerIdentity.Incognito,
				Agent:     strings.ToLower(player.CharacterID),
			}

			SetRank(p)

			playerRecv <- p
		}(player)
	}

	wg.Wait()
	close(playerRecv)

	// Populate players array
	// players := []*models.Player{}
	for player := range playerRecv {
		match.Players = append(match.Players, player)
	}

	SetPartyID(match.Players)
	SetNames(match.Players)

	SetPartySort(match.Players)

	return match
}

type GameFetchPlayerResp struct {
	MatchID string `json:"MatchID"`
}

type PregameFetchMatchResp struct {
	Teams []struct {
		Players []PlayerResp `json:"Players"`
	} `json:"Teams"`
	GamePodID string `json:"GamePodID"`
}

type PlayerResp struct {
	Subject                 string `json:"Subject"`
	TeamID                  string `json:"TeamID"`
	CharacterID             string `json:"CharacterID"`
	CharacterSelectionState string `json:"CharacterSelectionState"`
	PlayerIdentity          struct {
		AccountLevel     int  `json:"AccountLevel"`
		Incognito        bool `json:"Incognito"`
		HideAccountLevel bool `json:"HideAccountLevel"`
	} `json:"PlayerIdentity"`
}
