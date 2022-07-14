package riot

import (
	"net/http"
	"sync"

	"github.com/sovietscout/valorank/pkg/models"
)

type Pregame struct{
	GamePod string
}

func (r *Pregame) GetGamePod() string {
	return r.GamePod
}

func (r *Pregame) GetPlayers(playerChan chan <- []*models.Player) {
	// Current Pregame ID
	reqID, _ := http.NewRequest(http.MethodGet, GetGLZURL("/pregame/v1/players/" + UserPUUID), nil)
	reqID.Header = Local.GetRiotHeaders()

	respID, err := http.DefaultClient.Do(reqID)
	if err != nil {
		return
	}

	dataID := new(GameFetchPlayerResp)
	json.NewDecoder(respID.Body).Decode(dataID)

	// Pregame details
	req, _ := http.NewRequest(http.MethodGet, GetGLZURL("/pregame/v1/matches/" + dataID.MatchID), nil)
	req.Header = Local.GetRiotHeaders()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data := new(PregameFetchMatchResp)
	json.NewDecoder(resp.Body).Decode(data)

	var wg sync.WaitGroup
	playerRecv := make(chan *models.Player, len(data.Teams[0].Players))

	r.GamePod = data.GamePodID

	for _, player := range data.Teams[0].Players {
		wg.Add(1)

		go func(player PlayerResp) {
			defer wg.Done()

			p := &models.Player{
				SubjectID: player.Subject,
				Level: player.PlayerIdentity.AccountLevel,
				Ally: true,
				Incognito: player.PlayerIdentity.Incognito,
			}

			p.Agent = player.CharacterID

			SetRank(p)

			playerRecv <- p
		}(player)
	}

	wg.Wait()
	close(playerRecv)

	// Populate players array
	players := []*models.Player{}
	for player := range playerRecv {
		players = append(players, player)
	}

	SetNames(players)
	SetPartyID(players)
	SetLevelSort(players)

	playerChan <- players
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
	TeamID         			string `json:"TeamID"`
	CharacterID             string `json:"CharacterID"`
	CharacterSelectionState string `json:"CharacterSelectionState"`
	PlayerIdentity          struct {
		AccountLevel           int    `json:"AccountLevel"`
		Incognito              bool   `json:"Incognito"`
		HideAccountLevel       bool   `json:"HideAccountLevel"`
	} `json:"PlayerIdentity"`
}