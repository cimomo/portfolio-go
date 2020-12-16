package main

import "github.com/cimomo/portfolio-go/pkg/cmd"

func main() {
	portfolioCmd := cmd.NewPortfolioCmd()
	portfolioCmd.Execute()
}
