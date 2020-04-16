package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/cimomo/portfolio-go/pkg/portfolio"
	"github.com/cimomo/portfolio-go/pkg/terminal"
)

func main() {
	profile := flag.String("profile", "./examples/profile.yml", "(optional) Profile for portfolio")
	flag.Parse()

	p, err := loadPortfolio("Main", *profile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Hello, Portfolio", p.Holdings["VTI"].Status.Value)

	startTerminal(p)
}

func loadPortfolio(name string, profile string) (*portfolio.Portfolio, error) {
	p := portfolio.NewPortfolio("Main")

	err := p.Load(profile)
	if err != nil {
		return nil, err
	}

	p.Refresh()

	return p, nil
}

func startTerminal(portfolio *portfolio.Portfolio) {
	term := terminal.NewTerminal(portfolio)
	term.Start()
}
