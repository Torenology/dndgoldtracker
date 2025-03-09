package commands

import (
	"dndgoldtracker/models"
	"fmt"
)

// DistributeCoins prompts the user and splits coins among members
func DistributeCoins(p *models.Party) {
	if len(p.Members) == 0 {
		fmt.Println("No members to distribute coins to.")
		return
	}

	gold := GetIntUserInput("Gold: ")
	silver := GetIntUserInput("Silver: ")
	copper := GetIntUserInput("Copper: ")
	electrum := GetIntUserInput("Electrum: ")
	platinum := GetIntUserInput("Platinum: ")

	// Divide equally
	for i := range p.Members {
		p.Members[i].Wallet.Gold += gold / len(p.Members)
		p.Members[i].Wallet.Silver += silver / len(p.Members)
		p.Members[i].Wallet.Copper += copper / len(p.Members)
		p.Members[i].Wallet.Electrum += electrum / len(p.Members)
		p.Members[i].Wallet.Platinum += platinum / len(p.Members)
	}

	fmt.Println("Coins distributed!")
}

// AddXP distributes XP and checks for level-ups
func AddXP(p *models.Party) {
	if len(p.Members) == 0 {
		fmt.Println("No members to distribute XP to.")
		return
	}

	xp := GetIntUserInput("XP to distribute: ")
	for i := range p.Members {
		p.Members[i].XP += xp / len(p.Members)
		CheckLevelUp(&p.Members[i])
	}

	fmt.Println("XP added!")
}

// CheckLevelUp checks if a member levels up
func CheckLevelUp(member *models.Member) {
	var xpThresholds = []int{0, 300, 900, 2700, 6500, 14000, 23000, 34000, 48000, 64000, 85000, 100000, 120000, 140000, 165000, 195000, 225000, 265000, 305000, 355000} // XP values taken for D&D 5e

	for member.Level < len(xpThresholds)-1 {
		if member.XP >= xpThresholds[member.Level] {
			member.Level++
			fmt.Printf("🎉 %s leveled up to Level %d! 🎉\n", member.Name, member.Level)
		} else {
			break
		}
	}
}
