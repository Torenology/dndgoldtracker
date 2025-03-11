package main

import (
	"dndgoldtracker/ui"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Initialize and run the program
	p := tea.NewProgram(ui.NewModel())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
