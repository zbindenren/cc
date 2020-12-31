package main

import (
	"log"

	"github.com/zbindenren/cc/internal/cmd"
)

func main() {
	command := cmd.New()
	if err := command.Run(); err != nil {
		log.Fatal(err)
	}
}
