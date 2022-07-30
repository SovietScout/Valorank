package valorank

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sovietscout/valorank/pkg/content"
	"github.com/sovietscout/valorank/pkg/models"
	"golang.org/x/term"
)

var (
	Width, Height, _ = term.GetSize(int(os.Stdout.Fd()))

	faintStyle = lipgloss.NewStyle().Faint(true)
	centerStyle = lipgloss.NewStyle().Align(lipgloss.Center)

	gamePod string
	showRefresh = false
)

func (m *model) View() string {
	doc := strings.Builder{}

	doc.WriteString(m.infoLine() + "\n")
	doc.WriteString(m.tableLines())

	return doc.String()
}

func (m *model) infoLine() string {
	name := content.NAME + " " + content.VERSION + " |"
	help := m.helpLine()
	state := m.stateLine()

	sideW := max(name, help)

	line := lipgloss.JoinHorizontal(
		lipgloss.Center,
		faintStyle.Width(sideW).Align(lipgloss.Left).Render(name),
		centerStyle.Width(Width - (2 * sideW)).Render(state),
		faintStyle.Width(sideW).Align(lipgloss.Right).Render(help),
	)

	return line
}

func (m *model) stateLine() string {
	doc := strings.Builder{}

	doc.WriteString("Status: ")
	doc.WriteString(content.ColourFromState(m.client.State))

	switch m.client.State {
	case models.PREGAME:
	case models.INGAME:
		doc.WriteString(" (" + gamePod + ")")
	}

	return doc.String()
}

func (m *model) helpLine() string {
	hl := "| "

	if showRefresh {
		hl += "Fetching players…"
	} else {
		hl += m.help.View(m.keys)
	}

	return hl
}

func (m *model) tableLines() string {
	m.table.Table = m.table.Table.WithTargetWidth(Width)
	return m.table.View()
}

func (m *model) showRefresh(show bool) {
	showRefresh = show
}

func (m *model) setGamePodID(matchGPID string) {
	gamePod = content.ServerFromGamePod(matchGPID)
}

func max(a, b string) int {
	aW := lipgloss.Width(a)
	bW :=  lipgloss.Width(b)
	
	if aW >= bW {
		return aW
	} else {
		return bW
	}
}