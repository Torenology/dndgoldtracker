package commands

import (
	"dndgoldtracker/models"
	"testing"
)

// Test XP Distribution
func TestDistributeExperience(t *testing.T) {
	party := models.Party{
		Members: []models.Member{
			{Name: "Alice", Level: 1, XP: 0},
			{Name: "Bob", Level: 1, XP: 0},
		},
	}

	xpToAdd := 100
	expectedXP := xpToAdd / len(party.Members)

	DistributeExperience(&party, xpToAdd)

	// Check if XP was distributed correctly
	for _, member := range party.Members {
		if member.XP != expectedXP {
			t.Errorf("Expected %d XP, but got %d for %s", expectedXP, member.XP, member.Name)
		}
	}
}

func TestDistributeCoins(t *testing.T) {
	// Create a mock party with 3 members
	party := models.Party{
		Members: []models.Member{
			{Name: "Keg", CoinPriority: 0, Coins: make(map[string]int)},
			{Name: "Rowan", CoinPriority: 1, Coins: make(map[string]int)},
			{Name: "Fred", CoinPriority: 2, Coins: make(map[string]int)},
		},
	}

	// Coins to distribute
	money := map[string]int{
		platinum: 10,
		gold:     8,
		electrum: 7,
		silver:   5,
		copper:   3,
	}

	// Call the function to distribute the coins
	DistributeCoins(&party, money)

	// Check the distribution of coins
	tests := []struct {
		memberName       string
		expectedPlatinum int
		expectedGold     int
		expectedElectrum int
		expectedSilver   int
		expectedCopper   int
	}{
		{"Keg", 4, 3, 2, 2, 1},
		{"Rowan", 3, 2, 3, 2, 1},
		{"Fred", 3, 3, 2, 1, 1},
	}

	// Iterate through the test cases and compare expected vs actual
	for _, test := range tests {
		t.Run(test.memberName, func(t *testing.T) {
			member := getMemberByName(party.Members, test.memberName)

			if member.Coins[platinum] != test.expectedPlatinum {
				t.Errorf("%s's platinum: expected %d, got %d", test.memberName, test.expectedPlatinum, member.Coins[platinum])
			}
			if member.Coins[gold] != test.expectedGold {
				t.Errorf("%s's gold: expected %d, got %d", test.memberName, test.expectedGold, member.Coins[gold])
			}
			if member.Coins[electrum] != test.expectedElectrum {
				t.Errorf("%s's electrum: expected %d, got %d", test.memberName, test.expectedElectrum, member.Coins[electrum])
			}
			if member.Coins[silver] != test.expectedSilver {
				t.Errorf("%s's silver: expected %d, got %d", test.memberName, test.expectedSilver, member.Coins[silver])
			}
			if member.Coins[copper] != test.expectedCopper {
				t.Errorf("%s's copper: expected %d, got %d", test.memberName, test.expectedCopper, member.Coins[copper])
			}
		})
	}
}

// Helper function to get PartyMember by name
func getMemberByName(members []models.Member, name string) *models.Member {
	for i := range members {
		if members[i].Name == name {
			return &members[i]
		}
	}
	return nil
}
