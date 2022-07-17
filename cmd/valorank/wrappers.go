package valorank

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sovietscout/valorank/pkg/models"
)

func (m *model) waitForStateChange(ret chan models.State) tea.Cmd {
	if !m.client.ClientStateLoopOn {
		go m.client.ClientStateChangeLoop(ret)
	}

	return func() tea.Msg {
		return clientStateMsg(<-ret)
	}
}

func (m *model) waitForClientReady() tea.Cmd {
	ret := make(chan struct{})
	go m.client.ClientReadyLoop(ret)

	return func() tea.Msg {
		return clientReadyMsg(<-ret)
	}
}

func (m *model) waitForGameState(ret chan bool) tea.Cmd {
	return func() tea.Msg {
		return gameStateMsg(<-ret)
	}
}

func (m *model) waitForPlayers(ret chan []*models.Player) tea.Cmd {
	return func() tea.Msg {
		return <-ret
	}
}

type clientStateMsg models.State
type gameStateMsg bool
type clientReadyMsg struct{}
