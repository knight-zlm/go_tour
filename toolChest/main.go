package main

import (
	"log"

	"github.com/knight-zlm/go-tour/toolChest/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("cmd.Execute err: %v", err)
	}
}
