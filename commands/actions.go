package commands

import (
	"dndgoldtracker/models"
	"fmt"
	"sort"
)

const (
	platinum string = "Platinum"
	gold     string = "Gold"
	electrum string = "Electrum"
	silver   string = "Silver"
	copper   string = "Copper"
)

// DistributeCoins distributes coins fairly among party members in a fixed order
func DistributeCoins(p *models.Party, money map[string]int) {
	// Define the fixed order of coins
	coinOrder := []string{platinum, gold, electrum, silver, copper}

	numMembers := len(p.Members)
	if numMembers == 0 {
		fmt.Println("No members to distribute coins to.")
		return
	}

	// Initialize coin maps if not already set
	for i := range p.Members {
		if p.Members[i].Coins == nil {
			p.Members[i].Coins = make(map[string]int)
		}
	}

	// Helper function to distribute a specific coin type
	distributeCoin := func(coinType string, coinAmount int) {
		each := coinAmount / numMembers
		remainder := coinAmount % numMembers

		// Assign evenly to each member
		for i := range p.Members {
			fmt.Printf("Adding %d of %s to %s's wallet\n", each, coinType, p.Members[i].Name)
			p.Members[i].Coins[coinType] += each
		}

		// Sort members by priority for distributing the remainder
		sort.Slice(p.Members, func(i, j int) bool {
			return p.Members[i].CoinPriority < p.Members[j].CoinPriority
		})

		// Distribute excess coins based on priority
		for i := range remainder {
			p.Members[i].Coins[coinType]++
		}

		// Rotate priority to balance future distributions
		for i := range p.Members {
			p.Members[i].CoinPriority = (p.Members[i].CoinPriority + 1) % numMembers
		}
	}

	// Distribute coins in the predefined order
	for _, coinType := range coinOrder {
		amount, exists := money[coinType]
		if exists {
			distributeCoin(coinType, amount)
		}
	}
}

// DistributeExperience distributes XP and checks for level-ups
func DistributeExperience(p *models.Party, xp int) {
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
