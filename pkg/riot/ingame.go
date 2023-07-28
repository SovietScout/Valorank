package riot

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/sovietscout/valorank/pkg/local"
	"github.com/sovietscout/valorank/pkg/models"
)

type Ingame struct {
	sync.Mutex
}

func (r *Ingame) GetMatch() models.Match {
	r.Lock()
	defer r.Unlock()

	match := models.Match{State: models.INGAME}

	// Current Match ID
	reqID, _ := http.NewRequest(http.MethodGet, GetGLZURL("/core-game/v1/players/"+UserPUUID), nil)
	reqID.Header = local.Client.GetRiotHeaders()

	respID, err := client.Do(reqID)
	if err != nil {
		match.Err = err
		return match
	}
	defer respID.Body.Close()

	dataID := new(GameFetchPlayerResp)
	json.NewDecoder(respID.Body).Decode(dataID)

	// Match details
	req, _ := http.NewRequest(http.MethodGet, GetGLZURL("/core-game/v1/matches/"+dataID.MatchID), nil)
	req.Header = local.Client.GetRiotHeaders()

	resp, err := client.Do(req)
	if err != nil {
		match.Err = err
		return match
	}
	defer resp.Body.Close()

	data := new(CoregameFetchMatchResp)
	json.NewDecoder(resp.Body).Decode(data)

	match.GamePodID = data.GamePodID

	playerRecv := make(chan *models.Player, len(data.Players))
	var wg sync.WaitGroup

	for _, player := range data.Players {
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
				Agent:     strings.ToLower(player.CharacterID),
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
	// players := []*models.Player{}
	for player := range playerRecv {
		match.Players = append(match.Players, player)
	}

	SetPartyID(match.Players)
	SetNames(match.Players)

	SetPartySort(match.Players)
	SetTeamSort(match.Players)

	return match
}

type CoregameFetchMatchResp struct {
	Players   []PlayerResp `json:"Players"`
	GamePodID string       `json:"GamePodID"`
}
