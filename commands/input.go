package commands

import (
	"dndgoldtracker/models"
	"fmt"
	"strconv"
	"strings"
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
	newMember := models.Member{Name: name, Level: 1, XP: 0, CoinPriority: len(p.Members)}
	p.Members = append(p.Members, newMember)
	fmt.Println("Member added!")
}

// AddExperience allows the user to add xp to the created party.
func AddExperience(p *models.Party) {
	if len(p.Members) == 0 {
		fmt.Println("No members to distribute XP to.")
		return
	}

	xp := GetIntUserInput("XP to distribute: ")

	DistributeExperience(p, xp)
}

func AddCoins(p *models.Party) {
	if len(p.Members) == 0 {
		fmt.Println("No members to distribute coins to.")
		return
	}

	var money map[string]int = make(map[string]int)

	money[platinum] = GetIntUserInput("Platinum: ")
	money[gold] = GetIntUserInput("Gold: ")
	money[electrum] = GetIntUserInput("Electrum: ")
	money[silver] = GetIntUserInput("Silver: ")
	money[copper] = GetIntUserInput("Copper: ")

	DistributeCoins(p, money)
}

// Function to parse user input for coins
func ParseCoinsInput(input string) map[string]int {
	coins := make(map[string]int)
	coinParts := strings.Fields(input)
	for _, part := range coinParts {
		coinInfo := strings.Split(part, ":")
		if len(coinInfo) == 2 {
			coinType := coinInfo[0]
			amount, err := strconv.Atoi(coinInfo[1])
			if err == nil {
				coins[coinType] = amount
			}
		}
	}
	return coins
}
