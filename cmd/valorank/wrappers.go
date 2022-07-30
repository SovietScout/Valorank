package valorank

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sovietscout/valorank/pkg/models"
)

func (m *model) waitForStateChange(ret chan models.State) tea.Cmd {
	return func() tea.Msg {
		return clientStateMsg(<-ret)
	}
}

func (m *model) waitForApplicationState(ret chan bool) tea.Cmd {
	return func() tea.Msg {
		return applicationStateMsg(<-ret)
	}
}

func (m *model) waitForMatch(ret chan models.Match) tea.Cmd {
	return func() tea.Msg {
		return <-ret
	}
}

type clientStateMsg models.State
type applicationStateMsg bool
