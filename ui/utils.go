package ui

import (
	"dndgoldtracker/models"
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func checkbox(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}

// Convert members to table rows
func membersToRows(members []models.Member) []table.Row {
	var rows []table.Row
	for _, m := range members {
		rows = append(rows, table.Row{
			m.Name,
			strconv.Itoa(m.XP),
			strconv.Itoa(m.Coins[platinum]),
			strconv.Itoa(m.Coins[gold]),
			strconv.Itoa(m.Coins[electrum]),
			strconv.Itoa(m.Coins[silver]),
			strconv.Itoa(m.Coins[copper]),
		})
	}
	return rows
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func ConfigureInputs() []textinput.Model {

	i := make([]textinput.Model, 5)
	var t textinput.Model
	for j := range i {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch j {
		case 0:
			t.Placeholder = platinum
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = gold
		case 2:
			t.Placeholder = electrum
		case 3:
			t.Placeholder = silver
		case 4:
			t.Placeholder = copper
		}
		i[j] = t
	}

	return i
}

func ConfigureTable(members []models.Member) table.Model {
	columns := []table.Column{
		{Title: name, Width: 10},
		{Title: xp, Width: 6},
		{Title: platinum, Width: 10},
		{Title: gold, Width: 6},
		{Title: electrum, Width: 10},
		{Title: silver, Width: 8},
		{Title: copper, Width: 8},
	}

	rows := membersToRows(members)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(5),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t
}

func UpdateTableData(members []models.Member, t *table.Model) *table.Model {
	rows := membersToRows(members)
	t.SetRows(rows)
	return t
}

func ResetInputs(inputs []textinput.Model) {
	for i := range inputs {
		inputs[i].Reset()
	}
}
