package riot

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/sovietscout/valorank/pkg/models"
)

type Menu struct{}

func (r *Menu) GetGamePod() string {
	return ""
}

func (r *Menu) GetPlayers(playerChan chan <- []*models.Player) {
	userPartyID := getUserPartyID()

	req, _ := http.NewRequest(http.MethodGet, GetGLZURL("/parties/v1/parties/" + userPartyID), nil)
	req.Header = Local.GetRiotHeaders()

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data := new(FetchPartyResp)
	json.NewDecoder(resp.Body).Decode(data)

	var wg sync.WaitGroup
	playerRecv := make(chan *models.Player, len(data.Members))

	// Set players and their rank infos
	for _, player := range data.Members {
		wg.Add(1)

		go func(player MemberResp, partyID string) {
			defer wg.Done()

			p := &models.Player{
				SubjectID: player.Subject,
				PartyID: partyID,
				Level: player.PlayerIdentity.AccountLevel,
				Ally: true,
			}

			SetRank(p)

			playerRecv <- p
		}(player, userPartyID)
	}

	wg.Wait()
	close(playerRecv)

	// Populate players array
	players := []*models.Player{}
	for player := range playerRecv {
		players = append(players, player)
	}

	SetNames(players)
	SetLevelSort(players)

	playerChan <- players
}

func getUserPartyID() string {
	req, _ := http.NewRequest(http.MethodGet, GetGLZURL("/parties/v1/players/"+UserPUUID), nil)
	req.Header = Local.GetRiotHeaders()

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