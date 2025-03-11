package ui

import (
	"dndgoldtracker/models"
	"dndgoldtracker/storage"
	"fmt"
	"slices"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	platinum = "Platinum"
	gold     = "Gold"
	electrum = "Electrum"
	silver   = "Silver"
	copper   = "Copper"
	dotChar  = " • "
)

var (
	baseStyle     = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))
	subtleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	checkboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	dotStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
)

type model struct {
	table    table.Model
	party    models.Party
	Choice   int
	Chosen   bool
	Quitting bool
}

func (m model) Init() tea.Cmd { return nil }

// NewModel initializes the application state
func NewModel() model {
	party, err := storage.LoadParty() // Load saved data
	if err != nil {
		fmt.Println("Starting new party...")
		party = models.Party{}
	}

	members := party.Members
	columns := []table.Column{
		{Title: "Name", Width: 10},
		{Title: "XP", Width: 6},
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

	return model{
		party: party,
		table: t,
	}
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	if !m.Chosen {
		return updateChoices(msg, m)
	}

	switch m.Choice {
	case 0:
		return updateMoney(msg, m)
	case 1:
		return updateExperience(msg, m)
	case 2:
		m.Quitting = true
		return m, tea.Quit
	default:
		return m, nil
	}
}

// The main view, which just calls the appropriate sub-view
func (m model) View() string {
	var s string
	if m.Quitting {
		return "\n  See you later!\n\n"
	}

	if !m.Chosen {
		s = choicesView(m)
	} else {
		switch m.Choice {
		case 0:
			s = moneyView(m)
		case 1:
			s = xpView(m)
		case 2:
			s = "Bye, loser"
		default:
			s = "Don't do that"
		}
	}

	return baseStyle.Render("\n" + s + "\n\n")
}

// Sub-Updates

// Update loop for the first view where you're choosing a task.
func updateChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.Choice++
			if m.Choice > 2 {
				m.Choice = 0
			}
		case "k", "up":
			m.Choice--
			if m.Choice < 0 {
				m.Choice = 2
			}
		case "enter":
			m.Chosen = true
			return m, nil
		}
	}

	return m, nil
}

// Update loop for updating party money
func updateMoney(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	return m, nil
}

// Update loop for updating party experience
func updateExperience(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	return m, nil
}

// Sub-Views

// The first view, where you're choosing a task
func choicesView(m model) string {
	choice := m.Choice
	var msg string
	msg += baseStyle.Render(m.table.View())

	msg += "\nHave you earned money or experience?"
	msg += "\n"

	msg += fmt.Sprintf(
		"\n%s\n%s\n%s\n",
		checkbox("Money", choice == 0),
		checkbox("Experience", choice == 1),
		checkbox("Quit", choice == 2),
	)

	msg += subtleStyle.Render("j/k, up/down: select") + dotStyle +
		subtleStyle.Render("enter: choose") + dotStyle +
		subtleStyle.Render("q, esc: quit")

	return fmt.Sprintf(msg, nil)
}

func checkbox(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}

// The view for adding money
func moneyView(m model) string {
	var msg string
	currentPrio := slices.IndexFunc(m.party.Members, func(m models.Member) bool { return m.CoinPriority == 0 })
	msg += "Money entered here will be distributed to all party members as equally as possible.\n"
	msg += "Extra coins are distributed based on a priority system that rotates.\n"
	msg += fmt.Sprintf("Current Coin Priority is to %s\n\n", m.party.Members[currentPrio].Name)

	return msg
}

// The view for adding experience
func xpView(m model) string {
	var msg string
	msg += "Xp entered here will be distributed to all party members equally"

	return msg
}
