package playertable

import (
	"github.com/muesli/termenv"
	"github.com/sovietscout/valorank/pkg/models"
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
)

func GenerateParties(players []*models.Player) map[string]string {
	recvCh := make(chan bool, 1)
	iconCh := newIcon(recvCh)

	var tempPIDs = map[string]struct{}{}
	var dupPIDs = map[string]string{}

	// First pass: check which party IDs are duplicates (2+ players in the same party)
	for _, player := range players {
		if _, ok := tempPIDs[player.PartyID]; ok {
			recvCh <- true
			dupPIDs[player.PartyID] = <-iconCh
		}

		tempPIDs[player.PartyID] = struct{}{}
	}

	// Second pass: assign party icons if party ID is duplicate
	/*
		var partyIcons []string
		for _, player := range players {
			if icon, ok := dupPIDs[player.PartyID]; ok {
				partyIcons = append(partyIcons, icon)
			} else {
				partyIcons = append(partyIcons, "")
			}
		}
	*/

	recvCh <- false
	close(recvCh)

	return dupPIDs
}

func newIcon(recvCh <-chan bool) <-chan string {
	ch := make(chan string, 1)

	go func() {
		index := 0

		for send := range recvCh {
			if send {
				ch <- symbol.Foreground(partyIconColours[index]).String()
			} else {
				break
			}

			index++
		}

		close(ch)
	}()

	return ch
}
