package profile

import (
	"encoding/json"
	"testing"
)

func TestUserProfileJSON(t *testing.T) {
	jsonStr := `{"language": "English", "readKana": true, "levelFrom": "N5", "levelTo": "N4"}`
	var user UserProfile
	err := json.Unmarshal([]byte(jsonStr), &user)
	if err != nil {
		t.Errorf("Failed to unmarshal JSON: %v", err)
	}

	if user.Language != "English" {
		t.Errorf("Expected Language English, got %s", user.Language)
	}
	if !user.ReadKana {
		t.Errorf("Expected ReadKana true, got false")
	}
	if user.LevelFrom != "N5" {
		t.Errorf("Expected LevelFrom N5, got %s", user.LevelFrom)
	}
	if user.LevelTo != "N4" {
		t.Errorf("Expected LevelTo N4, got %s", user.LevelTo)
	}
}
