package main

import (
	"flag"
	"log"

	"github.com/cimomo/portfolio-go/pkg/portfolio"
	"github.com/cimomo/portfolio-go/pkg/terminal"
)

const (
	benchmark      = "SPY"
	initialBalance = 100000.00
)

func main() {
	profileFile := flag.String("profile", "./examples/profile.yml", "(optional) Profile for portfolio")
	flag.Parse()

	profile, err := loadProfile("Main", *profileFile)
	if err != nil {
		log.Fatal(err)
	}

	mkt := portfolio.NewMarket()

	perf := portfolio.NewPerformance(profile.Portfolios[0], benchmark, initialBalance)

	err = startTerminal(mkt, profile.Portfolios[0], perf)
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

func startTerminal(market *portfolio.Market, portfolio *portfolio.Portfolio, performance *portfolio.Performance) error {
	term := terminal.NewTerminal(market, portfolio, performance)

	err := term.Start()
	if err != nil {
		return err
	}

	return nil
}
