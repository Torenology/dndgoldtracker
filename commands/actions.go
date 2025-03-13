package commands

import (
	"dndgoldtracker/models"
	"fmt"
	"log"
	"slices"
	"sort"
)

// DistributeCoins distributes coins fairly among party members in a fixed order
func DistributeCoins(p *models.Party, money map[string]int) {
	numMembers := len(p.ActiveMembers)
	if numMembers == 0 {
		log.Println("No members to distribute coins to.")
		return
	}

	// Initialize coin maps if not already set
	for i := range p.ActiveMembers {
		if p.ActiveMembers[i].Coins == nil {
			p.ActiveMembers[i].Coins = make(map[string]int)
		}
	}

	// Helper function to distribute a specific coin type
	distributeCoin := func(coinType string, coinAmount int) {
		each := coinAmount / numMembers
		remainder := coinAmount % numMembers

		// Assign evenly to each member
		for i := range p.ActiveMembers {
			fmt.Printf("Adding %d of %s to %s's wallet\n", each, coinType, p.ActiveMembers[i].Name)
			p.ActiveMembers[i].Coins[coinType] += each
		}

		// Sort members by priority for distributing the remainder
		sort.Slice(p.ActiveMembers, func(i, j int) bool {
			return p.ActiveMembers[i].CoinPriority < p.ActiveMembers[j].CoinPriority
		})

		// Distribute excess coins based on priority
		for i := range remainder {
			p.ActiveMembers[i].Coins[coinType]++
		}

		// Rotate priority to balance future distributions
		for i := range p.ActiveMembers {
			p.ActiveMembers[i].CoinPriority = (p.ActiveMembers[i].CoinPriority + 1) % numMembers
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
	for i := range p.ActiveMembers {
		p.ActiveMembers[i].XP += xp / len(p.ActiveMembers)
		CheckLevelUp(&p.ActiveMembers[i])
	}

	fmt.Println("XP added!")
}

func AddMember(p *models.Party, name string, xp int, money map[string]int) {
	m := models.Member{Name: name, Level: DetermineLevel(xp), XP: xp, Coins: money, CoinPriority: len(p.ActiveMembers)}
	p.ActiveMembers = append(p.ActiveMembers, m)
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
			return i
		}
	}

	// max level
	return 20
}

func ChangeMemberActivity(p *models.Party, memberName string, tableName string) {
	log.Printf("tableName = %s, memberName = %s", tableName, memberName)
	if tableName == "Active" {
		// Deactivating Member
		memberIndex := slices.IndexFunc(p.ActiveMembers, func(m models.Member) bool { return m.Name == memberName })
		p.InactiveMembers = append(p.InactiveMembers, (p.ActiveMembers)[memberIndex])
		p.ActiveMembers = slices.Delete((p.ActiveMembers), memberIndex, memberIndex+1)
	} else {
		// Activating Member
		memberIndex := slices.IndexFunc(p.InactiveMembers, func(m models.Member) bool { return m.Name == memberName })
		p.ActiveMembers = append(p.ActiveMembers, (p.InactiveMembers)[memberIndex])
		p.InactiveMembers = slices.Delete((p.InactiveMembers), memberIndex, memberIndex+1)
	}
}
