package commands

import (
	"dndgoldtracker/models"
	"testing"
)

// Test XP Distribution
func TestDistributeExperience(t *testing.T) {
	party := models.Party{
		ActiveMembers: []models.Member{
			{Name: "Alice", Level: 1, XP: 0},
			{Name: "Bob", Level: 1, XP: 0},
		},
	}

	xpToAdd := 100
	expectedXP := xpToAdd / len(party.ActiveMembers)

	DistributeExperience(&party, xpToAdd)

	// Check if XP was distributed correctly
	for _, member := range party.ActiveMembers {
		if member.XP != expectedXP {
			t.Errorf("Expected %d XP, but got %d for %s", expectedXP, member.XP, member.Name)
		}
	}
}

func TestDistributeCoins(t *testing.T) {
	// Create a mock party with 3 members
	party := models.Party{
		ActiveMembers: []models.Member{
			{Name: "Keg", CoinPriority: 0, Coins: make(map[string]int)},
			{Name: "Rowan", CoinPriority: 1, Coins: make(map[string]int)},
			{Name: "Fred", CoinPriority: 2, Coins: make(map[string]int)},
		},
	}

	// Coins to distribute
	money := map[string]int{
		models.Platinum: 10,
		models.Gold:     8,
		models.Electrum: 7,
		models.Silver:   5,
		models.Copper:   3,
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
			member := getMemberByName(party.ActiveMembers, test.memberName)

			if member.Coins[models.Platinum] != test.expectedPlatinum {
				t.Errorf("%s's platinum: expected %d, got %d", test.memberName, test.expectedPlatinum, member.Coins[models.Platinum])
			}
			if member.Coins[models.Gold] != test.expectedGold {
				t.Errorf("%s's gold: expected %d, got %d", test.memberName, test.expectedGold, member.Coins[models.Gold])
			}
			if member.Coins[models.Electrum] != test.expectedElectrum {
				t.Errorf("%s's electrum: expected %d, got %d", test.memberName, test.expectedElectrum, member.Coins[models.Electrum])
			}
			if member.Coins[models.Silver] != test.expectedSilver {
				t.Errorf("%s's silver: expected %d, got %d", test.memberName, test.expectedSilver, member.Coins[models.Silver])
			}
			if member.Coins[models.Copper] != test.expectedCopper {
				t.Errorf("%s's copper: expected %d, got %d", test.memberName, test.expectedCopper, member.Coins[models.Copper])
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
