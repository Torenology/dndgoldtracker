package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Wallet holds different coin types
type Wallet struct {
	Gold     int
	Silver   int
	Copper   int
	Electrum int
	Platinum int
}

// Member represents a party member
type Member struct {
	Name   string
	Level  int
	XP     int
	Wallet Wallet
}

// Party represents the group of members
type Party struct {
	Members []Member
}

// SaveParty saves the party data to a JSON file
func (p *Party) SaveParty(filename string) error {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// LoadParty loads the party data from a JSON file
func LoadParty(filename string) (*Party, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var party Party
	err = json.Unmarshal(data, &party)
	if err != nil {
		return nil, err
	}
	return &party, nil
}

// DistributeCoins distributes coins evenly among party members
func (p *Party) DistributeCoins(coinType string, amount int) {
	count := len(p.Members)
	if count == 0 {
		fmt.Println("No members in the party!")
		return
	}

	each := amount / count
	remainder := amount % count

	for i := range p.Members {
		switch coinType {
		case "gold":
			p.Members[i].Wallet.Gold += each
		case "silver":
			p.Members[i].Wallet.Silver += each
		case "copper":
			p.Members[i].Wallet.Copper += each
		case "electrum":
			p.Members[i].Wallet.Electrum += each
		case "platinum":
			p.Members[i].Wallet.Platinum += each
		}
	}

	// Distribute remainder fairly
	for i := 0; i < remainder; i++ {
		switch coinType {
		case "gold":
			p.Members[i].Wallet.Gold++
		case "silver":
			p.Members[i].Wallet.Silver++
		case "copper":
			p.Members[i].Wallet.Copper++
		case "electrum":
			p.Members[i].Wallet.Electrum++
		case "platinum":
			p.Members[i].Wallet.Platinum++
		}
	}

	fmt.Printf("Distributed %d %s among party members.\n", amount, coinType)
}

// DistributeMultipleCoins allows distributing multiple coin types at once
func (p *Party) DistributeMultipleCoins(gold, silver, copper, electrum, platinum int) {
	count := len(p.Members)
	if count == 0 {
		fmt.Println("No members in the party!")
		return
	}

	// Function to distribute a single type of coin
	distribute := func(amount int, updateWallet func(*Wallet, int)) {
		each := amount / count
		remainder := amount % count

		// Give each member an equal share
		for i := range p.Members {
			updateWallet(&p.Members[i].Wallet, each)
		}

		// Distribute remainder fairly
		for i := 0; i < remainder; i++ {
			updateWallet(&p.Members[i].Wallet, 1)
		}
	}

	// Distribute each type of coin
	distribute(gold, func(w *Wallet, amount int) { w.Gold += amount })
	distribute(silver, func(w *Wallet, amount int) { w.Silver += amount })
	distribute(copper, func(w *Wallet, amount int) { w.Copper += amount })
	distribute(electrum, func(w *Wallet, amount int) { w.Electrum += amount })
	distribute(platinum, func(w *Wallet, amount int) { w.Platinum += amount })

	fmt.Println("Coins have been distributed among the party members.")
}

// XP thresholds for each level (D&D 5e example)
var xpThresholds = []int{0, 300, 900, 2700, 6500, 14000, 23000, 34000, 48000, 64000, 85000, 100000, 120000, 140000, 165000, 195000, 225000, 265000, 305000, 355000}

// AddExperience distributes XP among members and checks for level-ups
func (p *Party) AddExperience(amount int) {
	count := len(p.Members)
	if count == 0 {
		fmt.Println("No members in the party!")
		return
	}

	each := amount / count
	remainder := amount % count

	for i := range p.Members {
		p.Members[i].XP += each
		p.CheckLevelUp(&p.Members[i])
	}

	// Distribute remainder XP fairly
	for i := 0; i < remainder; i++ {
		p.Members[i].XP++
		p.CheckLevelUp(&p.Members[i])
	}

	fmt.Printf("Added %d XP to the party.\n", amount)
}

// CheckLevelUp checks if a member qualifies for a level-up
func (p *Party) CheckLevelUp(member *Member) {
	for member.Level < len(xpThresholds)-1 {
		if member.XP >= xpThresholds[member.Level] {
			member.Level++
			fmt.Printf("🎉 %s has leveled up to Level %d! 🎉\n", member.Name, member.Level)
		} else {
			break
		}
	}
}

// DisplayParty shows party details
func (p *Party) DisplayParty() {
	fmt.Println("\nParty Members:")
	for _, member := range p.Members {
		fmt.Printf("Name: %s | Level: %d | XP: %d | Wallet: %+v\n", member.Name, member.Level, member.XP, member.Wallet)
	}
	fmt.Println()
}

// AddMember allows the user to add a new party member
func (p *Party) AddMember(name string) {
	newMember := Member{Name: name, Level: 1, XP: 0, Wallet: Wallet{}}
	p.Members = append(p.Members, newMember)
	fmt.Printf("Added %s to the party!\n", name)
}

// GetUserInput prompts the user for input
func GetUserInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// GetIntUserInput prompts the user for an integer input
func GetIntUserInput(prompt string) int {
	input := GetUserInput(prompt)
	value, err := strconv.Atoi(input)
	if err != nil || value < 0 {
		fmt.Println("Invalid number, defaulting to 0.")
		return 0
	}
	return value
}

func main() {
	party := Party{}
	filename := "party.json"

	// Load party data if available
	if savedParty, err := LoadParty(filename); err == nil {
		party = *savedParty
		fmt.Println("Loaded existing party data.")
	} else {
		fmt.Println("No saved data found, starting fresh.")
	}

	for {
		fmt.Println("\nMenu:")
		fmt.Println("1. Add Party Member")
		fmt.Println("2. Distribute Coins")
		fmt.Println("3. Add Experience")
		fmt.Println("4. Show Party")
		fmt.Println("5. Save and Exit")
		choice := GetUserInput("Choose an option: ")

		switch choice {
		case "1":
			name := GetUserInput("Enter the member's name: ")
			party.AddMember(name)
		case "2":
			fmt.Println("Enter the amount of each coin type to distribute:")
			gold := GetIntUserInput("Gold: ")
			silver := GetIntUserInput("Silver: ")
			copper := GetIntUserInput("Copper: ")
			electrum := GetIntUserInput("Electrum: ")
			platinum := GetIntUserInput("Platinum: ")

			party.DistributeMultipleCoins(gold, silver, copper, electrum, platinum)
		case "3":
			amountStr := GetUserInput("Enter XP to distribute: ")
			amount, err := strconv.Atoi(amountStr)
			if err != nil || amount < 0 {
				fmt.Println("Invalid XP amount. Try again.")
				continue
			}
			party.AddExperience(amount)
		case "4":
			party.DisplayParty()
		case "5":
			err := party.SaveParty(filename)
			if err != nil {
				fmt.Println("Error saving party:", err)
			} else {
				fmt.Println("Party saved successfully!")
			}
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice. Try again.")
		}
	}
}
