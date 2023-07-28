package playertable

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/sovietscout/valorank/pkg/content"
	"github.com/sovietscout/valorank/pkg/models"
)

type Model struct {
	Table table.Model
}

const (
	cParty = "party"
	cAgent = "agent"
	cName  = "name"
	cRank  = "rank"
	cRR    = "rr"
	cPR    = "pr"
	cLevel = "level"
)

var (
	borderThin = table.Border{
		Top:    "─",
		Left:   "│",
		Right:  "│",
		Bottom: "─",

		TopRight:    "┐",
		TopLeft:     "┌",
		BottomRight: "┘",
		BottomLeft:  "└",

		TopJunction:    "┬",
		LeftJunction:   "├",
		RightJunction:  "┤",
		BottomJunction: "┴",
		InnerJunction:  "┼",

		InnerDivider: "│",
	}
)

func New() Model {
	centerStyle := lipgloss.NewStyle().Align(lipgloss.Center)

	return Model{
		Table: table.New([]table.Column{
			table.NewColumn(cParty, "Party", 7),
			table.NewFlexColumn(cAgent, "Agent", 1),
			table.NewFlexColumn(cName, "Name", 4),
			table.NewFlexColumn(cRank, "Rank", 1),
			table.NewColumn(cRR, "RR", 4),
			table.NewFlexColumn(cPR, "Peak Rank", 1),
			table.NewColumn(cLevel, "Level", 7),
		}).WithBaseStyle(centerStyle).HeaderStyle(centerStyle.Copy().Bold(true)).Border(borderThin),
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.Table, cmd = m.Table.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case models.Match:
		m.Table = m.Table.WithRows(generateRowsFromData(msg.Players))
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.Table.View()
}

func (m Model) Resize(width int) Model {
	m.Table = m.Table.WithTargetWidth(width)
	return m
}

func (m Model) Clear() Model {
	m.Table = m.Table.WithRows(generateRowsFromData([]*models.Player{}))
	return m
}

func generateRowsFromData(players []*models.Player) []table.Row {
	var playerParties = GenerateParties(players)

	var playerNameGen = NameGen()
	var lastPlayerTeam bool

	rows := []table.Row{}
	for i, player := range players {
		// Add empty row on team change
		if i != 0 && lastPlayerTeam != player.Ally {
			rows = append(rows, table.NewRow(table.RowData{}))
		}
		lastPlayerTeam = player.Ally

		// Player name based on incognito and whether agent has been selected
		agent := content.AgentFromID(player.Agent)
		if player.Incognito {
			if agent != nil {
				player.Name = playerNameGen(agent.Name)
			} else {
				player.Name = playerNameGen("Player")
			}
		}

		var partyIcon string
		if icon, ok := playerParties[player.PartyID]; ok {
			partyIcon = icon
		}

		row := table.NewRow(table.RowData{
			cParty: partyIcon,
			cAgent: content.StyleAgent(agent),
			cName:  ColouredName(player),
			cRank:  content.RankFromID(player.Rank),
			cRR:    player.RR,
			cPR:    content.RankFromID(player.PeakRank),
			cLevel: player.Level,
		})

		rows = append(rows, row)
	}

	return rows
}
