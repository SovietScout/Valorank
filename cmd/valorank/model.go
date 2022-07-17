package valorank

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sovietscout/valorank/pkg/client"
	"github.com/sovietscout/valorank/pkg/models"
	"github.com/sovietscout/valorank/pkg/models/components/playertable"
)

type model struct {
	client          *client.Client
	gameStateChan   chan (bool)
	clientStateChan chan (models.State)
	playerChan		chan ([]*models.Player)
	table			playertable.Model
	keys            KeyMap
	help            help.Model
}

func (m model) Init() tea.Cmd {
	go m.GameState()

	return tea.Batch(
		m.waitForGameState(m.gameStateChan),
	)
}

func NewModel() *model {
	playerChan := make(chan []*models.Player)

	return &model{
		client:          client.NewClient(playerChan),
		gameStateChan:   make(chan bool),
		clientStateChan: make(chan models.State),
		playerChan:		 playerChan,
		table:			 playertable.New(),
		keys:            DefaultKeyMap,
		help:            help.New(),
	}
}
