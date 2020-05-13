package main

import (
	"flag"
	"log"

	"github.com/cimomo/portfolio-go/pkg/portfolio"
	"github.com/cimomo/portfolio-go/pkg/terminal"
)

func main() {
	profile := flag.String("profile", "./examples/profile.yml", "(optional) Profile for portfolio")
	flag.Parse()

	port, err := loadPortfolio("Main", *profile)
	if err != nil {
		log.Fatal(err)
	}

	mkt := portfolio.NewMarket()

	perf := portfolio.NewPerformance(port)

	err = startTerminal(mkt, port, perf)
	if err != nil {
		log.Fatal(err)
	}
}

func loadPortfolio(name string, profile string) (*portfolio.Portfolio, error) {
	p := portfolio.NewPortfolio("Main")

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
