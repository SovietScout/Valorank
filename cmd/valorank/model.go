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
	applicationStateChan chan (bool)
	clientStateChan      chan (models.State)
	matchChan            chan (models.Match)
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
	matchChan := make(chan models.Match)

	return &model{
		client:               client.NewClient(matchChan),
		applicationStateChan: make(chan bool),
		clientStateChan:      make(chan models.State),
		matchChan:            matchChan,
		table:                playertable.New(),
		keys:                 DefaultKeyMap,
		help:                 help.New(),
	}
}
