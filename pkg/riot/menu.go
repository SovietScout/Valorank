package riot

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/sovietscout/valorank/pkg/local"
	"github.com/sovietscout/valorank/pkg/models"
)

type Menu struct{
	sync.Mutex
}

func (r *Menu) GetMatch() models.Match {
	r.Lock()
	defer r.Unlock()

	match := models.Match{State: models.MENU}

	UserPartyID := getUserPartyID()

	req, _ := http.NewRequest(http.MethodGet, GetGLZURL("/parties/v1/parties/"+UserPartyID), nil)
	req.Header = local.Client.GetRiotHeaders()

	resp, err := client.Do(req)
	if err != nil {
		match.Err = err
		return match
	}
	defer resp.Body.Close()

	data := new(FetchPartyResp)
	json.NewDecoder(resp.Body).Decode(data)

	// match.GamePodID = ""

	playerRecv := make(chan *models.Player, len(data.Members))
	var wg sync.WaitGroup

	// Set players and their rank infos
	for _, player := range data.Members {
		wg.Add(1)

		go func(player MemberResp, partyID string) {
			defer wg.Done()

			p := &models.Player{
				SubjectID: player.Subject,
				PartyID:   partyID,
				Level:     player.PlayerIdentity.AccountLevel,
				Ally:      true,
			}

			SetRank(p)

			playerRecv <- p
		}(player, UserPartyID)
	}

	wg.Wait()
	close(playerRecv)

	// Populate players array
	// players := []*models.Player{}
	for player := range playerRecv {
		match.Players = append(match.Players, player)
	}

	SetNames(match.Players)

	return match
}

func getUserPartyID() string {
	req, _ := http.NewRequest(http.MethodGet, GetGLZURL("/parties/v1/players/"+UserPUUID), nil)
	req.Header = local.Client.GetRiotHeaders()

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	data := new(FetchPlayerResp)
	json.NewDecoder(resp.Body).Decode(data)

	return data.CurrentPartyID
}

type FetchPartyResp struct {
	Members []MemberResp `json:"Members"`
}

type MemberResp struct {
	Subject         string `json:"Subject"`
	CompetitiveTier int    `json:"CompetitiveTier"`
	PlayerIdentity  struct {
		AccountLevel int  `json:"AccountLevel"`
		Incognito    bool `json:"Incognito"`
	} `json:"PlayerIdentity"`
}
