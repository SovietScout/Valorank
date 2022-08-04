package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
	"github.com/sovietscout/valorank/cmd/valorank"
	"github.com/sovietscout/valorank/pkg/content"
)

var (
	debug = flag.Bool("debug", true, "Enable debug mode")
)

func main() {
	flag.Parse()

	if *debug {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	termenv.SetWindowTitle(content.NAME + " by " + content.AUTHOR)

	p := tea.NewProgram(valorank.NewModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatalln("fatal: ", err)
	}
}