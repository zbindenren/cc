package main

import (
	"log"

	"github.com/zbindenren/cc/internal/cmd"
)

// nolint: gochecknoglobals
// following variables are set during build by goreleaser
var (
	version = "dev"
	commit  = "12345678"
	date    = "2006-01-02T15:04:05Z07:00"
)

func main() {
	b, err := cmd.NewBuildInfo(version, date, commit)
	if err != nil {
		log.Fatal(err)
	}

	command := cmd.New(*b)
	if err := command.Run(); err != nil {
		log.Fatal(err)
	}
}
