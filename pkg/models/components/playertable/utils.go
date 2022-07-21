package playertable

import (
	"github.com/muesli/termenv"
	"github.com/sovietscout/valorank/pkg/models"
)

var (
	symbol = termenv.String("■")
	p = termenv.ColorProfile()

	partyIconColours = []termenv.Color{
		p.Color("#E34343"),
		p.Color("#D843E3"),
		p.Color("#4346E3"),
		p.Color("#43E3D0"),
		p.Color("#5EE343"),
		p.Color("#E2ED39"),
		p.Color("#D452CF"),
	}
)

func GeneratePartyIcons(players []*models.Player) []string {
	stopCh := make(chan struct{})
	iconCh := newIcon(stopCh)

	var tempPIDs = map[string]struct{}{}
	var dupPIDs = map[string]string{}

	// First pass: check which party IDs are duplicates (2+ players in the same party)
	for _, player := range players {
		if _, ok := tempPIDs[player.PartyID]; ok {
			dupPIDs[player.PartyID] = <- iconCh
		}

		tempPIDs[player.PartyID] = struct{}{}
	}

	// Second pass: assign party icons if party ID is duplicate
	var partyIcons []string
	for _, player := range players {
		if icon, ok := dupPIDs[player.PartyID]; ok {
			partyIcons = append(partyIcons, icon)
		} else {
			partyIcons = append(partyIcons, "")
		}
	}

	stopCh <- struct{}{}
	close(stopCh)

	return partyIcons
}

func newIcon(stopCh <-chan struct{}) <- chan string {
	ch := make(chan string, 1)

	go func() {
		loop:
		for _, col := range partyIconColours {
			select {
			case <-stopCh:
				break loop
			default:
				ch <- symbol.Foreground(col).String()
			}
		}

		close(ch)
	}()

	return ch
}