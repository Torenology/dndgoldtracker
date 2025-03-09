package commands

import (
	"dndgoldtracker/models"
	"fmt"
	"strconv"
)

// GetUserChoice reads a string input
func GetUserChoice(prompt string) string {
	var input string
	fmt.Print(prompt)
	fmt.Scanln(&input)
	return input
}

// GetIntUserInput reads an integer safely
func GetIntUserInput(prompt string) int {
	fmt.Print(prompt)
	var input string
	fmt.Scanln(&input)

	value, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid input, defaulting to 0.")
		return 0
	}
	return value
}

// AddMember allows the user to create a new party member
func AddMember(p *models.Party) {
	name := GetUserChoice("Enter member name: ")
	newMember := models.Member{Name: name, Level: 1, XP: 0}
	p.Members = append(p.Members, newMember)
	fmt.Println("Member added!")
}
