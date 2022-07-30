package valorank

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"

	"github.com/sovietscout/valorank/pkg/models"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	// Keeps an eye on the gamestate (whether lockfile is present or not)
	case applicationStateMsg:
		if bool(msg) {
			if !m.client.IsRunning {
				m.client.Start(m.clientStateChan)
			}

			cmds = append(cmds,
				m.waitForApplicationState(m.applicationStateChan),
				m.waitForStateChange(m.clientStateChan),
			)
		} else {
			if m.client.IsRunning {
				m.client.Stop()
			}

			m.showRefresh(false)
			cmds = append(cmds, m.waitForApplicationState(m.applicationStateChan))
		}

	// Game state has changed
	case clientStateMsg:
		switch models.State(msg) {
		case models.OFFLINE:
			m.showRefresh(false)
			m.table = m.table.Clear()
		default:
			m.showRefresh(true)
			cmds = append(cmds, tea.Batch(
				m.waitForMatch(m.matchChan),
				m.waitForStateChange(m.clientStateChan),
			))
		}

	// Players received
	case models.Match:
		m.showRefresh(false)

		if msg.Err != nil {
			log.Println(msg.Err)
			break
		}

		m.setGamePodID(msg.GamePodID)

		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd, m.waitForMatch(m.matchChan))

	// Apparently not supported by Windows™
	case tea.WindowSizeMsg:
		width, _ := msg.Width, msg.Height
		m.table.Resize(width)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Quit):
			cmds = append(cmds, tea.Quit)
		case key.Matches(msg, DefaultKeyMap.Refresh):
			if m.client.State != models.OFFLINE {
				m.showRefresh(true)
				go m.client.GetMatch()
			}
		}

	}

	return m, tea.Batch(cmds...)
}

func (m *model) GameState() {
	var lockfileViable bool

	sendState := func(state bool) {
		lockfileViable = state
		m.applicationStateChan <- state
	}

	lockfilePath := filepath.Join(os.Getenv("LOCALAPPDATA"), "Riot Games/Riot Client/Config/lockfile")

	if _, err := os.Stat(lockfilePath); errors.Is(err, os.ErrNotExist) {
		sendState(false)
	} else {
		sendState(true)
	}

	watcher, _ := fsnotify.NewWatcher()

	done := make(chan bool)
	go func() {
		for event := range watcher.Events {
			switch {
			case event.Op&fsnotify.Write == fsnotify.Write:
				if event.Name == lockfilePath && !lockfileViable {
					sendState(true)
				}
			case event.Op&fsnotify.Remove == fsnotify.Remove:
				if event.Name == lockfilePath && lockfileViable {
					sendState(false)
				}
			}
		}
	}()

	watcher.Add(filepath.Dir(lockfilePath))

	<-done
	watcher.Close()
}
