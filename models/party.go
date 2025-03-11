package models

import "fmt"

const platinum string = "Platinum"
const gold string = "Gold"
const electrum string = "Electrum"
const silver string = "Silver"
const copper string = "Copper"

type Member struct {
	Name         string
	Level        int
	XP           int
	Coins        map[string]int
	CoinPriority int
}

type Party struct {
	Members []Member
}

// Display prints the current party state
func (p *Party) Display() {
	fmt.Println("\n=== Party Members ===")
	for _, member := range p.Members {
		fmt.Printf("%s (Level %d) - XP: %d, Wallet:%dPP %dGP %dEP %dSP %dCP \n",
			member.Name, member.Level, member.XP,
			member.Coins[platinum], member.Coins[gold], member.Coins[electrum],
			member.Coins[silver], member.Coins[copper])
	}
}
