package main

import (
	"fmt"

	"github.com/jonathanhecl/gollama"
)

var (
	model   = "phi4:latest"
	version = "0.1.0"
)

func main() {
	fmt.Println("Japanese Learning Agent v" + version)
	gl := gollama.New(model)
	if gl == nil {
		fmt.Println("Error creating Gollama instance")
		return
	}
	gl.Verbose = true

}
