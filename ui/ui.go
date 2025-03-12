package ui

import (
	"dndgoldtracker/commands"
	"dndgoldtracker/models"
	"dndgoldtracker/storage"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	name     = "Name"
	xp       = "XP"
	platinum = "Platinum"
	gold     = "Gold"
	electrum = "Electrum"
	silver   = "Silver"
	copper   = "Copper"
	dotChar  = " • "
)

var (
	baseStyle           = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))
	subtleStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	checkboxStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	dotStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	coins         = []string{platinum, gold, electrum, silver, copper}
)

type model struct {
	table          table.Model
	party          models.Party
	choice         int
	chosen         bool
	coinFocusIndex int
	coinInputs     []textinput.Model
	xpFocusIndex   int
	xpInputs       []textinput.Model
	cursorMode     cursor.Mode
	quitting       bool
}

func (m model) Init() tea.Cmd { return nil }

// NewModel initializes the application state
func NewModel() model {
	p, err := storage.LoadParty() // Load saved data
	if err != nil {
		fmt.Println("Starting new party...")
		p = models.Party{}
	}

	t := ConfigureTable(p.Members)
	ci := ConfigureCoinInputs()
	xi := ConfigureXpInputs()

	return model{
		party:      p,
		table:      t,
		coinInputs: ci,
		xpInputs:   xi,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	if !m.chosen {
		return updateChoices(msg, m)
	}

	switch m.choice {
	case 0:
		return updateMoney(msg, m)
	case 1:
		return updateExperience(msg, m)
	case 2:
		m.quitting = true
		return m, tea.Quit
	default:
		return m, nil
	}
}

// The main view, which just calls the appropriate sub-view
func (m model) View() string {
	var s string
	if m.quitting {
		return "\n  See you later!\n\n"
	}

	if !m.chosen {
		s = choicesView(m)
	} else {
		switch m.choice {
		case 0:
			s = moneyView(m)
		case 1:
			s = xpView(m)
		case 2:
			s = "Goodbye"
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
			m.choice++
			if m.choice > 2 {
				m.choice = 0
			}
		case "k", "up":
			m.choice--
			if m.choice < 0 {
				m.choice = 2
			}
		case "enter":
			m.chosen = true
			return m, nil
		}
	}

	return m, nil
}

// Update loop for updating party money
func updateMoney(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Change cursor mode
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.coinInputs))
			for i := range m.coinInputs {
				cmds[i] = m.coinInputs[i].Cursor.SetMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, Distribute money.
			if s == "enter" && m.coinFocusIndex == len(m.coinInputs) {
				var err error
				coinMap := make(map[string]int)
				// Set any unset values to 0
				for i := range m.coinInputs {
					if m.coinInputs[i].Value() == "" {
						m.coinInputs[i].SetValue("0")
					}
				}

				for i := range coins {
					coinMap[coins[i]], err = strconv.Atoi(m.coinInputs[i].Value())
					log.Printf("CoinMap entry for %s: %d\n", coins[i], coinMap[coins[i]])
					if err != nil {
						fmt.Printf("Invalid input for %s, try again\n", coins[i])
						return m, nil
					}
				}
				// Distribute the coins to the party
				commands.DistributeCoins(&m.party, coinMap)
				storage.SaveParty(&m.party)
				UpdateTableData(m.party.Members, &m.table)
				ResetInputs(m.coinInputs)

				m.chosen = false
				return m, nil
			}
			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.coinFocusIndex--
			} else {
				m.coinFocusIndex++
			}

			if m.coinFocusIndex > len(m.coinInputs) {
				m.coinFocusIndex = 0
			} else if m.coinFocusIndex < 0 {
				m.coinFocusIndex = len(m.coinInputs)
			}

			cmds := make([]tea.Cmd, len(m.coinInputs))
			for i := 0; i <= len(m.coinInputs)-1; i++ {
				if i == m.coinFocusIndex {
					// Set focused state
					cmds[i] = m.coinInputs[i].Focus()
					m.coinInputs[i].PromptStyle = focusedStyle
					m.coinInputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.coinInputs[i].Blur()
				m.coinInputs[i].PromptStyle = noStyle
				m.coinInputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}
	// Handle character input and blinking
	cmd := m.updateInputs(msg, m.coinInputs)

	return m, cmd
}

// Update loop for updating party experience
func updateExperience(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Change cursor mode
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.xpInputs))
			for i := range m.xpInputs {
				cmds[i] = m.xpInputs[i].Cursor.SetMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, Distribute xp.
			if s == "enter" && m.xpFocusIndex == len(m.xpInputs) {
				var err error
				if m.xpInputs[0].Value() == "" {
					m.xpInputs[0].SetValue("0")
				}

				xp, err := strconv.Atoi(m.xpInputs[0].Value())
				if err != nil {
					log.Println("Invalid input for experience, try again")
					return m, nil
				}

				commands.DistributeExperience(&m.party, xp)
				storage.SaveParty(&m.party)
				UpdateTableData(m.party.Members, &m.table)
				ResetInputs(m.xpInputs)

				m.chosen = false
				return m, nil
			}
			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.xpFocusIndex--
			} else {
				m.xpFocusIndex++
			}

			if m.xpFocusIndex > len(m.xpInputs) {
				m.xpFocusIndex = 0
			} else if m.xpFocusIndex < 0 {
				m.xpFocusIndex = len(m.xpInputs)
			}

			cmds := make([]tea.Cmd, len(m.xpInputs))
			for i := 0; i <= len(m.xpInputs)-1; i++ {
				if i == m.xpFocusIndex {
					// Set focused state
					cmds[i] = m.xpInputs[i].Focus()
					m.xpInputs[i].PromptStyle = focusedStyle
					m.xpInputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.xpInputs[i].Blur()
				m.xpInputs[i].PromptStyle = noStyle
				m.xpInputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}
	// Handle character input and blinking
	cmd := m.updateInputs(msg, m.xpInputs)

	return m, cmd
}

// Sub-Views

// The first view, where you're choosing a task
func choicesView(m model) string {
	choice := m.choice
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

	return msg
}

// The view for adding money
func moneyView(m model) string {
	var msg strings.Builder
	currentPrio := slices.IndexFunc(m.party.Members, func(m models.Member) bool { return m.CoinPriority == 0 })
	msg.WriteString("Money entered here will be distributed to all party members as equally as possible.")
	msg.WriteRune('\n')
	msg.WriteString("Extra coins are distributed based on a priority system that rotates.")
	msg.WriteRune('\n')
	msg.WriteString(fmt.Sprintf("Current Coin Priority is to %s\n\n", m.party.Members[currentPrio].Name))

	for i := range m.coinInputs {
		msg.WriteString(m.coinInputs[i].View())
		if i < len(m.coinInputs)-1 {
			msg.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.coinFocusIndex == len(m.coinInputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&msg, "\n\n%s\n\n", *button)

	msg.WriteString(helpStyle.Render("cursor mode is "))
	msg.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	msg.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return msg.String()
}

// The view for adding experience
func xpView(m model) string {
	var msg strings.Builder
	msg.WriteString("Xp entered here will be distributed to all party members equally")
	for i := range m.xpInputs {
		msg.WriteString(m.xpInputs[i].View())
		if i < len(m.xpInputs)-1 {
			msg.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.coinFocusIndex == len(m.coinInputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&msg, "\n\n%s\n\n", *button)

	msg.WriteString(helpStyle.Render("cursor mode is "))
	msg.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	msg.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return msg.String()
}
