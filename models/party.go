package models

import "fmt"

type Wallet struct {
	Gold, Silver, Copper, Electrum, Platinum int
}

type Member struct {
	Name   string
	Level  int
	XP     int
	Wallet Wallet
}

type Party struct {
	Members []Member
}

// Display prints the current party state
func (p *Party) Display() {
	fmt.Println("\n=== Party Members ===")
	for _, member := range p.Members {
		fmt.Printf("%s (Level %d) - XP: %d, Wallet: %dGP %dSP %dCP %dEP %dPP\n",
			member.Name, member.Level, member.XP,
			member.Wallet.Gold, member.Wallet.Silver, member.Wallet.Copper,
			member.Wallet.Electrum, member.Wallet.Platinum)
	}
}
