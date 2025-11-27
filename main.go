package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/jonathanhecl/gollama"
)

var (
	model   = "mistral:latest"
	version = "0.1.4"

	// User
	language  = ""
	readKana  = false
	levelFrom = ""
	levelTo   = ""
)

func main() {
	fmt.Println("Japanese Learning Agent v" + version)
	if v := os.Getenv("OLLAMA_MODEL"); v != "" {
		model = v
	}
	fmt.Println("Model: " + model)
	gl := gollama.New(model)
	if gl == nil {
		fmt.Println("Error creating Gollama instance")
		return
	}
	// gl.Verbose = true

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	gl.SetSystemPrompt("You are a helpful and patient Japanese language teacher. Your goal is to help the user learn Japanese based on their level and interests. Be encouraging and clear.")

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
	r, err = gl.Chat(ctx, "Analyze the user's response and extract their profile.\n"+
		"1. Native Language\n"+
		"2. Knowledge of Hiragana/Katakana\n"+
		"3. Current Level\n"+
		"4. Target Level\n\n"+
		"User response: \""+input+"\"\n\n"+
		"Instructions:\n"+
		"- Default LevelFrom: none\n"+
		"- Default LevelTo: JLPT N5\n"+
		"- Respond STRICTLY with valid JSON only. No markdown, no explanations.\n"+
		"- Use the following schema:\n"+
		"{\n"+
		"  \"language\": \"string\",\n"+
		"  \"readKana\": boolean,\n"+
		"  \"levelFrom\": \"string\",\n"+
		"  \"levelTo\": \"string\"\n"+
		"}", gollama.StructToStructuredFormat(userStr{}))
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}
	if err := r.DecodeContent(&user); err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	fmt.Println("Your language: " + user.Language)
	fmt.Println("You read kana: " + fmt.Sprintf("%t", user.ReadKana))
	fmt.Println("Your level from: " + user.LevelFrom)
	fmt.Println("Your level to: " + user.LevelTo)

	userPrompt := "[ROLE]\n"
	userPrompt += "You are a Japanese teacher. Teach the user based on their profile.\n\n"
	userPrompt += "[USER PROFILE]\n"
	userPrompt += "- Language: " + user.Language + "\n"
	userPrompt += "- Current Level: " + user.LevelFrom + "\n"
	userPrompt += "- Target Level: " + user.LevelTo + "\n"
	if user.ReadKana {
		userPrompt += "- Kana Knowledge: Yes (Use Kana, NO Romaji)\n"
	} else {
		userPrompt += "- Kana Knowledge: No (Use Romaji for all Japanese text)\n"
	}
	userPrompt += "\n[INSTRUCTIONS]\n"
	userPrompt += "1. Explain everything in " + user.Language + ".\n"
	userPrompt += "2. Do not use Kanji above " + user.LevelTo + ".\n"
	userPrompt += "3. If you use Kanji, always provide the reading.\n"
	userPrompt += "4. When providing a Japanese sentence, break it down grammatically.\n"
	userPrompt += "5. Be concise but helpful.\n\n"
	userPrompt += "[RESPONSE FORMAT]\n"
	userPrompt += "1. Japanese Sentence (with reading if needed)\n"
	userPrompt += "2. Breakdown (Bullet points for each word/particle)\n"
	userPrompt += "3. Translation\n"
	userPrompt += "4. Explanation (if requested or necessary)\n\n"
	userPrompt += "[EXAMPLE]\n"
	userPrompt += "Sentence: 今日は学校に行きます (Kyō wa gakkō ni ikimasu)\n"
	userPrompt += "Breakdown:\n"
	userPrompt += "* 今日 (Kyō) - Today [Noun]\n"
	userPrompt += "* は (wa) - Topic Marker [Particle]\n"
	userPrompt += "* 学校 (gakkō) - School [Noun]\n"
	userPrompt += "* に (ni) - Direction Marker [Particle]\n"
	userPrompt += "* 行きます (ikimasu) - To go [Verb, Polite]\n"
	userPrompt += "Translation: I am going to school today.\n"

	gl.SetSystemPrompt(userPrompt)

	fmt.Println("Welcome to the Japanese Learning Agent!")

	for {
		if ctx.Err() != nil {
			fmt.Println("\nさようなら")
			break
		}
		fmt.Print("You: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("\nさようなら")
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
			fmt.Println("\nさようなら")
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
