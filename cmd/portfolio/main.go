package main

import (
	"flag"
	"log"

	"github.com/cimomo/portfolio-go/pkg/portfolio"
	"github.com/cimomo/portfolio-go/pkg/terminal"
)

func main() {
	profileFile := flag.String("profile", "./examples/profile.yml", "(optional) Profile for portfolio")
	flag.Parse()

	profile, err := loadProfile("Main", *profileFile)
	if err != nil {
		log.Fatal(err)
	}

	err = startTerminal(profile)
	if err != nil {
		log.Fatal(err)
	}
}

func loadProfile(name string, profile string) (*portfolio.Profile, error) {
	p := portfolio.NewProfile(name)

	err := p.Load(profile)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func startTerminal(profile *portfolio.Profile) error {
	term := terminal.NewTerminal(profile)

	err := term.Start()
	if err != nil {
		return err
	}

	return nil
}
