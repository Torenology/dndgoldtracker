package ui

import (
	"dndgoldtracker/commands"
	"dndgoldtracker/models"
	"dndgoldtracker/storage"
	"fmt"
	"log"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

// Sub-Updates

// Update loop for the first view where you're choosing a task.
func updateChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.choice++
			if m.choice > 3 {
				m.choice = 0
			}
		case "k", "up":
			m.choice--
			if m.choice < 0 {
				m.choice = 3
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

func updateActivateMembers(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var selectedTable table.Model
	var tableName string
	if m.activeMemberTable.Focused() {
		selectedTable = m.activeMemberTable
		tableName = "Active"
	} else {
		selectedTable = m.inactiveMemberTable
		tableName = "Inactive"
	}
	var activeCmd tea.Cmd
	var inactiveCmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if m.activeMemberTable.Focused() {
				m.activeMemberTable.Blur()
				m.inactiveMemberTable.Focus()
			} else {
				m.activeMemberTable.Focus()
				m.inactiveMemberTable.Blur()
			}
		case "enter":
			// activate/deactivate member
			selectedMemberName := selectedTable.SelectedRow()[0]
			commands.ChangeMemberActivity(&m.party, selectedMemberName, tableName)
			m.activeMemberTable.SetRows(membersToRows(m.party.ActiveMembers))
			m.inactiveMemberTable.SetRows(membersToRows(m.party.InactiveMembers))
		case "s":
			storage.SaveParty(&m.party)
			m.chosen = false
		}

	}
	m.activeMemberTable, activeCmd = m.activeMemberTable.Update(msg)
	m.inactiveMemberTable, inactiveCmd = m.inactiveMemberTable.Update(msg)
	return m, tea.Batch(activeCmd, inactiveCmd)
}
