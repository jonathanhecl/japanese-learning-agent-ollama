package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jonathanhecl/gollama"
)

var (
	model   = "mistral:latest"
	version = "0.1.1"

	// User
	language  = ""
	readKana  = false
	levelFrom = ""
	levelTo   = ""
)

func main() {
	fmt.Println("Japanese Learning Agent v" + version)
	gl := gollama.New(model)
	if gl == nil {
		fmt.Println("Error creating Gollama instance")
		return
	}
	// gl.Verbose = true

	ctx := context.Background()
	gl.SetSystemPrompt("You are a Japanese teacher, you are teaching Japanese to a student.")

	type userStr struct {
		Language  string `json:"language" required:"true" description:"Language the user speaks"`
		ReadKana  bool   `json:"readKana" required:"true" description:"True if the user can read hiragana and katakana"`
		LevelFrom string `json:"levelFrom" required:"true" description:"none, N1, N2, N3, N4, N5, JLPT N1, JLPT N2, JLPT N3, JLPT N4, JLPT N5, etc"`
		LevelTo   string `json:"levelTo" required:"true" description:"N1, N2, N3, N4, N5, JLPT N1, JLPT N2, JLPT N3, JLPT N4, JLPT N5, etc"`
	}
	user := userStr{}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Thinking...")
	r, err := gl.Chat(ctx, "Ask the user what language they speak, if they understand kanas (hiragana and katakana), "+
		"and what their level of Japanese is (Noken, JLPT, etc) and what level the user wants to learn in order to tailor "+
		"the learning experience to the user. Ask IN ENGLISH!")
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}
	fmt.Println("Teacher: ", r.Content)
	fmt.Print("You: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			fmt.Println("さようなら")
			return
		}
		fmt.Println("Error reading input:", err)
	}
	input = strings.TrimSpace(input)

	fmt.Println("Thinking...")
	r, err = gl.Chat(ctx, "1. What is your native language or the primary language you speak?\n"+
		"2. Do you understand hiragana and katakana (the two kana systems in Japanese)?\n"+
		"3. Could you tell me about your current level of Japanese?\n"+
		"4. What level would you like to achieve in Japanese?\n\n"+

		"User response: "+input+"\n\n"+

		"By default, if not specified, LevelFrom is none.\n"+
		"By default, if not specified, LevelTo is JLPT N5.\n"+
		"Complete the user's profile with the response.\n"+
		"Answer in JSON format.", gollama.StructToStructuredFormat(userStr{}))
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}
	if err := r.DecodeContent(&user); err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	userPrompt := "The user speaks " + user.Language + ", so whenever you explain something to them, it must be in the user's language.\n"
	userPrompt += "The user has indicated that they are at level " + user.LevelFrom + " and wants to reach level " + user.LevelTo + ".\n"
	userPrompt += "Do not use kanji characters that are above the level the user wishes to learn. Instead, use kana.\n"
	userPrompt += "When you respond in Japanese, break down the sentence by explaining each element.\n"
	userPrompt += "If it contains kanji from level " + user.LevelTo + ", explain how it is read in kana. If there are particles, explain them.\n"
	userPrompt += "\nExample:\n今日は学校で友達と話しましたよ\n"
	userPrompt += "* [NOUN] 今日 (きょう, Meaning: Today)\n"
	userPrompt += "* [PARTICLE] は (Particle: Topic marker)\n"
	userPrompt += "* [NOUN] 学校 (がっこう, Meaning: School)\n"
	userPrompt += "* [PARTICLE] で (Particle: Location marker)\n"
	userPrompt += "* [NOUN] 友達 (ゆうだい, Meaning: Friend)\n"
	userPrompt += "* [PARTICLE] と (Particle: Connector)\n"
	userPrompt += "* [VERB] 話しました (はなしゃみした, Meaning: To speak) - 話 (はな) is the verb to speak. The verb is in past tense (ます form).\n"
	userPrompt += "* [PARTICLE] よ (Particle: Emphasis marker)\n"
	userPrompt += "Translation: Today I spoke with my friend at school.\n\n"
	if !user.ReadKana {
		userPrompt += "The user does not know kana, so you should use romaji instead.\n"
	} else {
		userPrompt += "The user can read kana. Then don't use romaji.\n"
	}

	gl.SetSystemPrompt("You are a Japanese teacher, you are teaching Japanese to a student.\n" + userPrompt)

	fmt.Println("Welcome to the Japanese Learning Agent!")

	for {
		fmt.Print("You: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("\nBye!")
				break
			}
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			// ignore empty inputs
			continue
		}
		if strings.EqualFold(input, "exit") || strings.EqualFold(input, "quit") {
			fmt.Println("さようなら")
			break
		}

		fmt.Println("Thinking...")
		res, err := gl.Chat(ctx, input)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			continue
		}
		fmt.Println("Teacher: " + res.Content)
	}

}
