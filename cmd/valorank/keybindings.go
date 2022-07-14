package valorank

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Quit key.Binding
	Refresh key.Binding
}

var DefaultKeyMap = KeyMap {
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl-c", "Quit Valorank"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("f5", "r"),
		key.WithHelp("F5", "Refresh players"),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Refresh}
}

// FullHelp returns keybindings for the expanded help view
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Refresh, k.Quit},
	}
}