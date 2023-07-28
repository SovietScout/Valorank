package valorank

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/sovietscout/valorank/pkg/content"
	"github.com/sovietscout/valorank/pkg/models"
	"golang.org/x/term"
)

var (
	Width, Height, _ = term.GetSize(int(os.Stdout.Fd()))

	faintStyle = lipgloss.NewStyle().Faint(true)
	centerStyle = lipgloss.NewStyle().Align(lipgloss.Center)

	gamePod string
)

func (m *model) View() string {
	doc := strings.Builder{}

	doc.WriteString(m.infoLine() + "\n")
	doc.WriteString(m.tableLines() + "\n")
	doc.WriteString(m.logLine())

	return doc.String()
}

func (m *model) infoLine() string {
	name := content.NAME + " " + content.VERSION
	help := m.help.View(m.keys)
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
	case models.PREGAME, models.INGAME:
		doc.WriteString(" (" + gamePod + ")")
	}

	return doc.String()
}

func (m *model) tableLines() string {
	m.table.Table = m.table.Table.WithTargetWidth(Width)
	return m.table.View()
}

func (m *model) logLine() string {
	if m.content != "" {
		log := centerStyle.Width(Width - 2).Render(m.content)
		return termenv.String("└" + log + "┘").Faint().String()
	}

	return m.content
}

func (m *model) setGamePodID(matchGPID string) {
	gamePod = content.ServerFromGamePod(matchGPID)
}

func max(a, b string) int {
	aW := lipgloss.Width(a)
	bW := lipgloss.Width(b)
	
	if aW >= bW {
		return aW
	}

	return bW
}