package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
	"github.com/sovietscout/valorank/cmd/valorank"
	"github.com/sovietscout/valorank/pkg/content"
)

func main() {
	termenv.SetWindowTitle(content.NAME + " by " + content.AUTHOR)

	p := tea.NewProgram(valorank.NewModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatalln("There has been an error: ", err)
	}
}