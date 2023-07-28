package ui

import (
	"github.com/charmbracelet/bubbles/help"
	// "github.com/sovietscout/valorank/pkg/types"
)

type App struct {
	keys help.Model

	// view map[types.Viewable]view.View
}

func New() *App {
	return &App{}
}