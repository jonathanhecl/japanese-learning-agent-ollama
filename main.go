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
	model    = "phi4:latest"
	language = "spanish"
	version  = "0.1.0"
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
	gl.SetSystemPrompt("You are a Japanese teacher, you are teaching Japanese to a student. The student speaks " + language + ".")

	reader := bufio.NewReader(os.Stdin)
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
			fmt.Println("Bye!")
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
