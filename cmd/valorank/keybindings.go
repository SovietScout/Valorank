package valorank

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Quit key.Binding
	Refresh key.Binding
	Help key.Binding
}

var DefaultKeyMap = KeyMap {
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "Quit Valorank"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r", "f5"),
		key.WithHelp("r", "Refresh players"),
	),
	Help: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "Help"),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Refresh}
}

// FullHelp returns keybindings for the expanded help view
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help, k.Refresh, k.Quit},
	}
}