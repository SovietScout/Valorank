package riot

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/sovietscout/valorank/pkg/models"
)

type Ingame struct {
	GamePod string
}

func (r *Ingame) GetGamePod() string {
	return r.GamePod
}

func (r *Ingame) GetPlayers(playerChan chan<- []*models.Player) {
	// Current Match ID
	reqID, _ := http.NewRequest(http.MethodGet, GetGLZURL("/core-game/v1/players/"+UserPUUID), nil)
	reqID.Header = Local.GetRiotHeaders()

	respID, err := http.DefaultClient.Do(reqID)
	if err != nil {
		return
	}
	defer respID.Body.Close()

	dataID := new(GameFetchPlayerResp)
	json.NewDecoder(respID.Body).Decode(dataID)

	// Match details
	req, _ := http.NewRequest(http.MethodGet, GetGLZURL("/core-game/v1/matches/"+dataID.MatchID), nil)
	req.Header = Local.GetRiotHeaders()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data := new(CoregameFetchMatchResp)
	json.NewDecoder(resp.Body).Decode(data)

	var wg sync.WaitGroup
	playerRecv := make(chan *models.Player, len(data.Players))

	r.GamePod = data.GamePodID

	for _, player := range data.Players {
		wg.Add(1)

		go func(player PlayerResp) {
			defer wg.Done()

			p := &models.Player{
				SubjectID: player.Subject,
				Level:     player.PlayerIdentity.AccountLevel,
				Agent:     player.CharacterID,
				Ally:      player.TeamID == "Blue",
				Incognito: player.PlayerIdentity.Incognito,
			}

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
	SetTeamSort(players)

	playerChan <- players
}

type CoregameFetchMatchResp struct {
	Players   []PlayerResp `json:"Players"`
	GamePodID string       `json:"GamePodID"`
}
