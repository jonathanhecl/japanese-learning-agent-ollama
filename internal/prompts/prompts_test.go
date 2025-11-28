package prompts

import (
	"strings"
	"testing"

	"github.com/jonathanhecl/japanese-learning-agent-ollama/internal/profile"
)

func TestConstructSystemPrompt(t *testing.T) {
	user := profile.UserProfile{
		Language:  "English",
		ReadKana:  true,
		LevelFrom: "N5",
		LevelTo:   "N4",
	}

	prompt := ConstructSystemPrompt(user)

	if !strings.Contains(prompt, "Language: English") {
		t.Error("Prompt should contain user language")
	}
	if !strings.Contains(prompt, "Kana Knowledge: Yes") {
		t.Error("Prompt should indicate Kana knowledge")
	}
	if !strings.Contains(prompt, "Target Level: N4") {
		t.Error("Prompt should contain target level")
	}
}

func TestGetProfileExtractionPrompt(t *testing.T) {
	input := "I speak Spanish"
	prompt := GetProfileExtractionPrompt(input)

	if !strings.Contains(prompt, "User response: \"I speak Spanish\"") {
		t.Error("Prompt should contain user input")
	}
	if !strings.Contains(prompt, "Respond STRICTLY with valid JSON only") {
		t.Error("Prompt should contain strict JSON instruction")
	}
}
