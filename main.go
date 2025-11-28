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

	"github.com/jonathanhecl/japanese-learning-agent-ollama/internal/agent"
	"github.com/jonathanhecl/japanese-learning-agent-ollama/internal/prompts"
)

var (
	model   = "mistral:latest"
	version = "0.1.4"
)

func main() {
	fmt.Println("Japanese Learning Agent v" + version)
	if v := os.Getenv("OLLAMA_MODEL"); v != "" {
		model = v
	}
	fmt.Println("Model: " + model)

	ag, err := agent.New(model)
	if err != nil {
		fmt.Println("Error creating agent:", err)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ag.SetSystemPrompt("You are a helpful and patient Japanese language teacher. Your goal is to help the user learn Japanese based on their level and interests. Be encouraging and clear.")

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Thinking...")
	// Initial greeting and profile gathering
	initialPrompt := "Ask the user what language they speak, if they understand kanas (hiragana and katakana), " +
		"and what their level of Japanese is (Noken, JLPT, etc) and what level the user wants to learn in order to tailor " +
		"the learning experience to the user. Ask IN ENGLISH!"

	r, err := ag.Chat(ctx, initialPrompt)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}
	fmt.Println("Teacher: ", r)
	fmt.Print("You: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			fmt.Println("さようなら")
			return
		}
		fmt.Println("Error reading input:", err)
		return
	}
	input = strings.TrimSpace(input)

	fmt.Println("Thinking...")
	userProfile, err := ag.ExtractProfile(ctx, input)
	if err != nil {
		fmt.Println("Error extracting profile: " + err.Error())
		return
	}

	fmt.Println("Your language: " + userProfile.Language)
	fmt.Println("You read kana: " + fmt.Sprintf("%t", userProfile.ReadKana))
	fmt.Println("Your level from: " + userProfile.LevelFrom)
	fmt.Println("Your level to: " + userProfile.LevelTo)

	systemPrompt := prompts.ConstructSystemPrompt(*userProfile)
	ag.SetSystemPrompt(systemPrompt)

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
			continue
		}
		if strings.EqualFold(input, "exit") || strings.EqualFold(input, "quit") {
			fmt.Println("\nさようなら")
			break
		}

		fmt.Println("Thinking...")
		res, err := ag.Chat(ctx, input)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			continue
		}
		fmt.Println("Teacher: " + res)
	}
}
