package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
	"github.com/sovietscout/valorank/cmd/valorank"
	"github.com/sovietscout/valorank/pkg/content"
)

var (
	debug = flag.Bool("debug", false, "Enable debug mode")
)

func main() {
	flag.Parse()

	log.SetOutput(io.Discard)

	if *debug {
		logsFolder := filepath.Join(".", "logs")

		if err := os.MkdirAll(logsFolder, os.ModePerm); err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}

		f, err := tea.LogToFile(filepath.Join(logsFolder, time.Now().Format("2006-01-02") + ".log"), "debug")

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