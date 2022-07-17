package valorank

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
	"github.com/sovietscout/valorank/pkg/client"
	"github.com/sovietscout/valorank/pkg/models"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	// Keeps an eye on the gamestate (whether lockfile is present or not)
	case gameStateMsg:
		if bool(msg) {
			cmds = append(cmds, tea.Batch(
				m.waitForGameState(m.gameStateChan),
				m.waitForClientReady(),
			))
		} else {
			// [Game is/has been turned off]
			// Remake client, throw all running functions into the garbage
			m.showRefresh(false)
			m.client = client.NewClient(m.playerChan)

			cmds = append(cmds, m.waitForGameState(m.gameStateChan))
		}


	// Game is ready to accept requests
	case clientReadyMsg:
		cmds = append(cmds, m.waitForStateChange(m.clientStateChan))


	// Game state has changed
	case clientStateMsg:
		if client.State(msg) == client.OFFLINE {
			m.showRefresh(false)
			m.table = m.table.Clear()
		} else {
			m.showRefresh(true)
			cmds = append(cmds, tea.Batch(
				m.waitForPlayers(m.playerChan),
				m.waitForStateChange(m.clientStateChan),
			))
		}


	// Players received
	case []*models.Player:
		m.showRefresh(false)
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, tea.Batch(
			cmd,
			m.waitForPlayers(m.playerChan),
		))


	// Apparently not supported by Windows™
	case tea.WindowSizeMsg:
		width, _ := msg.Width, msg.Height
		m.table.Resize(width)


	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Quit):
			cmds = append(cmds, tea.Quit)
		case key.Matches(msg, DefaultKeyMap.Refresh):
			if m.client.State != client.OFFLINE {
				m.showRefresh(true)
				go m.client.GetPlayers()
			}
		}

	}

	return m, tea.Batch(cmds...)
}

func (m *model) GameState() {
	lockfile_path := filepath.Join(os.Getenv("LOCALAPPDATA"), "Riot Games/Riot Client/Config/lockfile")

	if _, err := os.Stat(lockfile_path); errors.Is(err, os.ErrNotExist) {
		m.gameStateChan <- false
	} else {
		m.gameStateChan <- true
	}

	watcher, _ := fsnotify.NewWatcher()

	done := make(chan bool)
	go func() {
		for event := range watcher.Events {
			switch {
			case event.Op&fsnotify.Create == fsnotify.Create:
				if event.Name == lockfile_path {
					m.gameStateChan <- true
				}
			case event.Op&fsnotify.Remove == fsnotify.Remove:
				if event.Name == lockfile_path {
					m.gameStateChan <- false
				}
			}
		}
	}()

	watcher.Add(filepath.Dir(lockfile_path))

	<-done
	watcher.Close()
}
