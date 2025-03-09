package main

import (
	"dndgoldtracker/commands"
	"dndgoldtracker/models"
	"dndgoldtracker/storage"
	"fmt"
)

func main() {
	party, err := storage.LoadParty() // Load saved data
	if err != nil {
		fmt.Println("Starting new party...")
		party = models.Party{}
	}

	for {
		fmt.Println("\n=== D&D Party Manager ===")
		fmt.Println("1. View Party")
		fmt.Println("2. Add Member")
		fmt.Println("3. Distribute Coins")
		fmt.Println("4. Add XP")
		fmt.Println("5. Save and Exit")

		choice := commands.GetUserChoice("Choose an option: ")

		switch choice {
		case "1":
			party.Display()
		case "2":
			commands.AddMember(&party)
		case "3":
			commands.DistributeCoins(&party)
		case "4":
			commands.AddXP(&party)
		case "5":
			storage.SaveParty(party)
			fmt.Println("Party saved. Exiting...")
			return
		default:
			fmt.Println("Invalid choice, try again.")
		}
	}
}
