package valorank

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sovietscout/valorank/pkg/content"
	"golang.org/x/term"
)

var (
	Width, Height, _ = term.GetSize(int(os.Stdout.Fd()))

	faintStyle = lipgloss.NewStyle().Faint(true)
	centerStyle = lipgloss.NewStyle().Align(lipgloss.Center)

	showRefresh = false
)

func (m *model) View() string {
	doc := strings.Builder{}

	doc.WriteString(m.infoLine() + "\n")
	doc.WriteString(m.tableLines())

	return doc.String()
}

func (m *model) infoLine() (view string) {
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

func (m *model) stateLine() (view string) {
	doc := strings.Builder{}

	doc.WriteString("Status: ")
	doc.WriteString(content.ColourFromState(string(m.client.State)))

	if m.client.Riot != nil {
		if server := content.ServerFromGamePod(m.client.Riot.GetGamePod()); server != "" {
			doc.WriteString(" (")
			doc.WriteString(server)
			doc.WriteString(")")
		}
 	}

	return doc.String()
}

func (m *model) helpLine() (view string) {
	hl := "| "

	if showRefresh {
		hl += "Fetching players…"
	} else {
		hl += m.help.View(m.keys)
	}

	return hl
}

func (m *model) tableLines() (view string) {
	m.table.Table = m.table.Table.WithTargetWidth(Width)
	return m.table.View()
}

func (m *model) showRefresh(show bool) {
	showRefresh = show
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