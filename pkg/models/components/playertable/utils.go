package playertable

import (
	"strconv"

	"github.com/muesli/termenv"
	"github.com/sovietscout/valorank/pkg/models"
	"github.com/sovietscout/valorank/pkg/riot"
)

var (
	symbol = termenv.String("■")
	p      = termenv.ColorProfile()

	partyIconColours = []termenv.Color{
		p.Color("#E34343"),
		p.Color("#D843E3"),
		p.Color("#4346E3"),
		p.Color("#43E3D0"),
		p.Color("#5EE343"),
		p.Color("#E2ED39"),
		p.Color("#D452CF"),
	}

	teamColour = p.Color("#4c97ed")
	oppColour  = p.Color("#ee4d4d")
	userColour = p.Color("#dde029")
)

func iconGen() func() string {
	var index int = 0

	return func() string {
		defer func() { index++ }()
		return symbol.Foreground(partyIconColours[index]).String()
	}
}

func GenerateParties(players []*models.Player) map[string]string {
	newIcon := iconGen()

	var tempPIDs = map[string]struct{}{}
	var dupPIDs = map[string]string{}

	// First pass: check which party IDs are duplicates (2+ players in the same party)
	for _, player := range players {
		if _, ok := tempPIDs[player.PartyID]; ok {
			if dupPIDs[player.PartyID] == "" {
				dupPIDs[player.PartyID] = newIcon()
			}
		} else {
			tempPIDs[player.PartyID] = struct{}{}
		}
	}

	newIcon = nil
	return dupPIDs
}

func ColouredName(player *models.Player) string {
	var colour termenv.Color

	if player.SubjectID == riot.UserPUUID {
		colour = userColour
	} else if player.Ally {
		colour = teamColour
	} else {
		colour = oppColour
	}

	return termenv.String(player.Name).Foreground(colour).String()
}

func NameGen() func(string) string {
	var nameMap = map[string]int{}

	return func(s string) string {
		if _, ok := nameMap[s]; !ok {
			nameMap[s] = 1
			return s
		}

		defer func() { nameMap[s]++ }()
		return s + " " + strconv.Itoa(nameMap[s])
	}
}
