package commands

import (
	"dndgoldtracker/models"
	"fmt"
	"log"
	"sort"
)

// DistributeCoins distributes coins fairly among party members in a fixed order
func DistributeCoins(p *models.Party, money map[string]int) {
	numMembers := len(p.Members)
	if numMembers == 0 {
		log.Println("No members to distribute coins to.")
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
	for _, coinType := range models.CoinOrder {
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

func AddMember(p *models.Party, name string, xp int, money map[string]int) {
	m := models.Member{Name: name, Level: DetermineLevel(xp), XP: xp, Coins: money, CoinPriority: len(p.Members)}
	p.Members = append(p.Members, m)
	log.Printf("Welcome to the party %s!\n", m.Name)
}

// CheckLevelUp checks if a member levels up
func CheckLevelUp(member *models.Member) {

	for member.Level < len(models.XpThresholds)-1 {
		if member.XP >= models.XpThresholds[member.Level] {
			member.Level++
			log.Printf("🎉 %s leveled up to Level %d! 🎉\n", member.Name, member.Level)
		} else {
			break
		}
	}
}

func DetermineLevel(xp int) int {
	for i := range models.XpThresholds {
		if xp < models.XpThresholds[i] {
			return i + 1
		}
	}

	// max level
	return 20
}
