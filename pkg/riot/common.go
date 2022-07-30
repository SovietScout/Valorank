package riot

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/sovietscout/valorank/pkg/content"
	"github.com/sovietscout/valorank/pkg/models"
)

var (
	UserPUUID string
	Region    string

	Local *NetCL

	client = http.Client{Timeout: 10 * time.Second}
)

type RiotClient interface {
	GetMatch() models.Match
}

func SetRank(player *models.Player) error {
	req, _ := http.NewRequest(http.MethodGet, GetPDURL("/mmr/v1/players/"+player.SubjectID), nil)
	req.Header = Local.GetRiotHeaders()

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return errors.New("rate limited")
	}

	data := new(MMRFetchPlayerResp)
	json.NewDecoder(resp.Body).Decode(data)

	for seasonID, details := range data.QueueSkills.Competitive.SeasonalInfoBySeasonID {
		rankInSeason := 0

		if details.WinsByTier != nil {

			// Assign highest rank in that season
			for rankStr := range details.WinsByTier {
				if rankStr == "0" {
					continue
				}

				rank, _ := strconv.Atoi(rankStr)

				if rank > rankInSeason {
					rankInSeason = rank
				}
			}

			// Pre-ascendant modification
			if rankInSeason > 20 && content.PreAscendantSeasonID(seasonID) {
				rankInSeason += 3
			}

			if player.PeakRank < rankInSeason {
				player.PeakRank = rankInSeason
			}

		}

		if seasonID == content.CurrentSeasonID {
			player.Rank = details.CompetitiveTier
			player.RR = details.RankedRating
		}
	}

	return nil
}

func SetPartyID(players []*models.Player) error {
	resp, err := Local.GET("/chat/v4/presences")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data := new(PresencesResp)
	json.NewDecoder(resp.Body).Decode(data)

	for _, presence := range data.Presences {
		for _, player := range players {
			if presence.Puuid == player.SubjectID {
				private_bytes, _ := base64.StdEncoding.DecodeString(presence.Private)

				data := new(PresencesPrivate)
				json.Unmarshal(private_bytes, data)

				player.PartyID = data.PartyID
				break
			}
		}
	}

	return nil
}

// Get names (and taglines) of people who aren't incognito
func SetNames(players []*models.Player) error {
	PUUIDs := []string{}

	for _, player := range players {
		if !player.Incognito {
			PUUIDs = append(PUUIDs, player.SubjectID)
		}
	}

	jsonPUUIDs, _ := json.Marshal(PUUIDs)

	req, _ := http.NewRequest(http.MethodPut, GetPDURL("/name-service/v2/players"), bytes.NewBuffer(jsonPUUIDs))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var data []NameResp
	json.NewDecoder(resp.Body).Decode(&data)

	for i := range data {
		for _, player := range players {
			if data[i].Subject == player.SubjectID {
				playerNameData := data[i]
				player.Name = playerNameData.GameName + "#" + playerNameData.TagLine

				break
			}
		}
	}

	return nil
}

// Sorts players by level in descending order
func SetLevelSort(players []*models.Player) {
	sort.Slice(players, func(i, j int) bool {
		return players[i].Level > players[j].Level
	})
}

// Sorts players by team
func SetTeamSort(players []*models.Player) {
	sort.Slice(players, func(i, j int) bool {
		return players[i].Ally
	})
}

type FetchPlayerResp struct {
	CurrentPartyID string `json:"CurrentPartyID"`
}

type NameResp struct {
	Subject  string `json:"Subject"`
	GameName string `json:"GameName"`
	TagLine  string `json:"TagLine"`
}

type MMRFetchPlayerResp struct {
	QueueSkills struct {
		Competitive struct {
			SeasonalInfoBySeasonID map[string]SeasonIDResp `json:"SeasonalInfoBySeasonID"`
		} `json:"competitive"`
	} `json:"QueueSkills"`
}

type SeasonIDResp struct {
	SeasonID                   string         `json:"SeasonID"`
	NumberOfWins               int            `json:"NumberOfWins"`
	NumberOfWinsWithPlacements int            `json:"NumberOfWinsWithPlacements"`
	NumberOfGames              int            `json:"NumberOfGames"`
	Rank                       int            `json:"Rank"`
	CapstoneWins               int            `json:"CapstoneWins"`
	LeaderboardRank            int            `json:"LeaderboardRank"`
	CompetitiveTier            int            `json:"CompetitiveTier"`
	RankedRating               int            `json:"RankedRating"`
	WinsByTier                 map[string]int `json:"WinsByTier"`
	GamesNeededForRating       int            `json:"GamesNeededForRating"`
	TotalWinsNeededForRank     int            `json:"TotalWinsNeededForRank"`
}

type PresencesResp struct {
	Presences []struct {
		Private string `json:"private"`
		Product string `json:"product"`
		Puuid   string `json:"puuid"`
	} `json:"presences"`
}

type PresencesPrivate struct {
	SessionLoopState string `json:"sessionLoopState"`
	PartyID          string `json:"partyId"`
}
