package main

import (
	"flag"
	"log"

	"github.com/cimomo/portfolio-go/pkg/terminal"
)

func main() {
	profile := flag.String("profile", "./examples/profile.yml", "(optional) Profile for portfolio")
	flag.Parse()

	err := startTerminal(*profile)
	if err != nil {
		log.Fatal(err)
	}
}

func startTerminal(profile string) error {
	term := terminal.NewTerminal(profile)

	err := term.Start()
	if err != nil {
		return err
	}

	return nil
}
