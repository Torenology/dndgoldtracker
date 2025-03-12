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
	name    = "Name"
	xp      = "XP"
	level   = "Level"
	dotChar = " • "
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

	focusedButton   = focusedStyle.Render("[ Submit ]")
	blurredButton   = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	xpFields        = []string{xp}
	newMemberFields = []string{name, xp}
)

type model struct {
	table            table.Model
	party            models.Party
	choice           int
	chosen           bool
	coinFocusIndex   int
	coinInputs       []textinput.Model
	xpFocusIndex     int
	xpInputs         []textinput.Model
	memberFocusIndex int
	memberInputs     []textinput.Model
	cursorMode       cursor.Mode
	quitting         bool
}

func (m model) Init() tea.Cmd { return nil }

// NewModel initializes the application state
func NewModel() model {
	p, err := storage.LoadParty() // Load saved data
	if err != nil {
		fmt.Println("Starting new party...")
		p = models.Party{}
	}

	newMemberFields = append(newMemberFields, models.CoinOrder...)

	t := configureTable(p.Members)

	ci := configureInputs(models.CoinOrder)
	xi := configureInputs(xpFields)
	mi := configureInputs(newMemberFields)

	return model{
		party:        p,
		table:        t,
		coinInputs:   ci,
		xpInputs:     xi,
		memberInputs: mi,
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
		return updateAddMember(msg, m)
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
			s = addMemberView(m)
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
			var cmds []tea.Cmd
			cmds = changeCursorMode(m.xpInputs, &m.cursorMode)

			return m, tea.Batch(cmds...)
		// Set focus to next input
		case "enter":
			// Did the user press enter while the submit button was focused?
			// If so, Distribute money.
			if m.coinFocusIndex == len(m.coinInputs) {
				var err error
				coinMap := make(map[string]int)
				// Set any unset values to 0
				handleUnsetInputs(m.coinInputs)

				for i := range models.CoinOrder {
					coinMap[models.CoinOrder[i]], err = strconv.Atoi(m.coinInputs[i].Value())
					log.Printf("CoinMap entry for %s: %d\n", models.CoinOrder[i], coinMap[models.CoinOrder[i]])
					if err != nil {
						fmt.Printf("Invalid input for %s, try again\n", models.CoinOrder[i])
						return m, nil
					}
				}

				// Distribute the coins to the party
				commands.DistributeCoins(&m.party, coinMap)
				saveUpdateReset(&m)

				m.chosen = false
				return m, nil
			}
			// Cycle indexes
		case "up", "shift-tab", "down":
			s := msg.String()
			if s == "down" {
				m.coinFocusIndex++
			} else {
				m.coinFocusIndex--
			}
			log.Printf("coinFocusIndex = %d", m.coinFocusIndex)
			cmds := updateFocusIndex(&m.coinFocusIndex, m.coinInputs)
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
			var cmds []tea.Cmd
			cmds = changeCursorMode(m.xpInputs, &m.cursorMode)
			return m, tea.Batch(cmds...)

		// Set focus to next input
		case "enter":
			// Did the user press enter while the submit button was focused?
			// If so, Distribute xp.
			if m.xpFocusIndex == len(m.xpInputs) {
				var err error
				handleUnsetInputs(m.xpInputs)

				xp, err := strconv.Atoi(m.xpInputs[0].Value())
				if err != nil {
					log.Println("Invalid input for experience, try again")
					return m, nil
				}

				commands.DistributeExperience(&m.party, xp)
				saveUpdateReset(&m)

				m.chosen = false
				return m, nil
			}
		case "up", "shift-tab", "down":
			s := msg.String()
			if s == "down" {
				m.xpFocusIndex++
			} else {
				m.xpFocusIndex--
			}
			cmds := updateFocusIndex(&m.xpFocusIndex, m.xpInputs)
			return m, tea.Batch(cmds...)
		}
	}
	// Handle character input and blinking
	cmd := m.updateInputs(msg, m.xpInputs)

	return m, cmd
}

func updateAddMember(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Change cursor mode
		case "ctrl+r":
			var cmds []tea.Cmd
			cmds = changeCursorMode(m.memberInputs, &m.cursorMode)

			return m, tea.Batch(cmds...)
		// Set focus to next input
		case "enter":
			// Did the user press enter while the submit button was focused?
			// If so, Distribute money.
			if m.memberFocusIndex == len(m.memberInputs) {
				var err error
				if m.memberInputs[0].Value() == "" {
					log.Println("Name value required, try again")
					return m, nil
				}
				name := m.memberInputs[0].Value()
				// Set any unset values other than name to 0
				handleUnsetInputs(m.memberInputs)
				xp, err := strconv.Atoi(m.memberInputs[1].Value())
				if err != nil {
					log.Println("Error occurred with xp, try again")
					return m, nil
				}

				newMemberCoins := m.memberInputs[2:len(m.memberInputs)]
				newMemberMoney := make(map[string]int)
				for i := range newMemberCoins {
					newMemberMoney[newMemberCoins[i].Placeholder], err = strconv.Atoi(newMemberCoins[i].Value())
					if err != nil {
						log.Printf("Error occurred with %s. Please try again", newMemberCoins[i].Placeholder)
						return m, nil
					}
				}

				// Add the new party Member
				commands.AddMember(&m.party, name, xp, newMemberMoney)
				saveUpdateReset(&m)

				m.chosen = false
				return m, nil
			}
			// Cycle indexes
		case "up", "shift-tab", "down":
			s := msg.String()
			if s == "down" {
				m.memberFocusIndex++
			} else {
				m.memberFocusIndex--
			}
			log.Printf("memberFocusIndex = %d", m.memberFocusIndex)
			cmds := updateFocusIndex(&m.memberFocusIndex, m.memberInputs)
			return m, tea.Batch(cmds...)
		}
	}
	// Handle character input and blinking
	cmd := m.updateInputs(msg, m.memberInputs)

	return m, cmd
}

// Sub-Views

// The first view, where you're choosing a task
func choicesView(m model) string {
	choice := m.choice
	var msg string
	msg += baseStyle.Render(m.table.View())

	msg += "\nWhat would you like to do?"
	msg += "\n"

	msg += fmt.Sprintf(
		"\n%s\n%s\n%s\n",
		checkbox("Distribute Money", choice == 0),
		checkbox("Distribute Experience", choice == 1),
		checkbox("Add Member", choice == 2),
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
	msg.WriteString("Money entered here will be distributed to all party members as equally as possible.\n")
	msg.WriteString("Extra coins are distributed based on a priority system that rotates.\n")
	msg.WriteString(fmt.Sprintf("Current Coin Priority is to %s\n\n", m.party.Members[currentPrio].Name))
	msg.WriteString(buildInputList(m.coinInputs, m.coinFocusIndex, m.cursorMode))
	return msg.String()
}

// The view for adding experience
func xpView(m model) string {
	var msg strings.Builder
	msg.WriteString("Xp entered here will be distributed to all party members equally\n\n")
	msg.WriteString(buildInputList(m.xpInputs, m.xpFocusIndex, m.cursorMode))
	return msg.String()
}

func addMemberView(m model) string {
	var msg strings.Builder
	msg.WriteString("Enter the new party member's data\n")
	msg.WriteString(buildInputList(m.memberInputs, m.memberFocusIndex, m.cursorMode))
	return msg.String()
}
