package valorank

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sovietscout/valorank/pkg/client"
	"github.com/sovietscout/valorank/pkg/models"
	"github.com/sovietscout/valorank/pkg/models/components/playertable"
)

type model struct {
	client               *client.Client
	stateCh              chan (models.State)
	matchCh              chan (models.Match)
	applicationStateChan chan (bool)
	content              string
	table                playertable.Model
	keys                 KeyMap
	help                 help.Model
}

func (m model) Init() tea.Cmd {
	go m.GameState()

	return tea.Batch(
		m.waitForApplicationState(m.applicationStateChan),
	)
}

func NewModel() *model {
	matchCh := make(chan models.Match)
	stateCh := make(chan models.State)

	return &model{
		client:               client.NewClient(stateCh, matchCh),
		stateCh:              stateCh,
		matchCh:              matchCh,
		applicationStateChan: make(chan bool),
		content:              initializingStr,
		table:                playertable.New(),
		keys:                 DefaultKeyMap,
		help:                 help.New(),
	}
}
