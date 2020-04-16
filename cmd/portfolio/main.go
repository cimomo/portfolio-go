package main

import (
	"fmt"

	"github.com/cimomo/portfolio-go/pkg/portfolio"
	"github.com/cimomo/portfolio-go/pkg/terminal"
)

func main() {
	p := loadPortfolio("Main")

	fmt.Println("Hello, Portfolio", p.Holdings["VTI"].Status.Value)

	startTerminal(p)
}

func loadPortfolio(name string) *portfolio.Portfolio {
	p := portfolio.NewPortfolio("Main")
	p.Load("./examples/profile.yml")
	p.Refresh()

	return p
}

func startTerminal(portfolio *portfolio.Portfolio) {
	term := terminal.NewTerminal(portfolio)
	term.Start()
}
