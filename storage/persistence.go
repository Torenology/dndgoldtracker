package storage

import (
	"dndgoldtracker/models"
	"encoding/json"
	"os"
)

// SaveParty writes party data to a JSON file
func SaveParty(party *models.Party) error {
	data, err := json.MarshalIndent(party, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("party.json", data, 0644)
}

// LoadParty loads party data from a JSON file
func LoadParty() (models.Party, error) {
	data, err := os.ReadFile("party.json")
	if err != nil {
		return models.Party{}, err
	}
	var party models.Party
	err = json.Unmarshal(data, &party)
	return party, err
}
