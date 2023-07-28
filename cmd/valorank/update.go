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
	"github.com/sovietscout/valorank/pkg/riot"
)

const (
	fetchPlayerStr  = "Fetching players"
	initializingStr = "Initializing Valorank"
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
				m.client.Start()
			}

			cmds = append(cmds, m.waitForStateChange(m.stateCh),
				m.waitForApplicationState(m.applicationStateChan))
		} else {
			if m.client.IsRunning {
				m.client.Stop()
			}

			cmds = append(cmds, m.waitForApplicationState(m.applicationStateChan))
		}

		m.content = ""

	// Game state has changed
	case clientStateMsg:
		switch models.State(msg) {
		case models.OFFLINE:
			m.content = ""
			m.table = m.table.Clear()
		default:
			m.content = fetchPlayerStr
			cmds = append(cmds, tea.Batch(
				m.waitForMatch(m.matchCh),
				m.waitForStateChange(m.stateCh),
			))
		}

	// Players received
	case models.Match:
		m.content = ""

		if msg.Err != nil {
			log.Println(msg.Err)

			if errors.Is(msg.Err, riot.ErrPregameIndexOutOfRange) {
				m.content = "You have been rate limited"
			} else {
				m.content = "An error has occured. Please file an issue with the logs should the problem persist"
			}

			break
		}

		m.setGamePodID(msg.GamePodID)

		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd, m.waitForMatch(m.matchCh))

	/* Apparently not supported by Windows™
	case tea.WindowSizeMsg:
		width, _ := msg.Width, msg.Height
		m.table.Resize(width)
	*/

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Quit):
			log.Println("valorank: Program exited")
			cmds = append(cmds, tea.Quit)
		case key.Matches(msg, DefaultKeyMap.Refresh):
			if m.client.State != models.OFFLINE {
				log.Println("valorank: Manually refreshed")
				m.content = fetchPlayerStr
				go m.client.GetMatch()
			}
		}

	}

	return m, tea.Batch(cmds...)
}

func (m *model) GameState() {
	var lockfileViable bool

	sendState := func(state bool) {
		log.Printf("valorank: Lockfile present: %t\n", state)

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

func Init() {
	
}
