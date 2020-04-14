package main

import (
	"fmt"

	"github.com/cimomo/portfolio-go/pkg/portfolio"
)

func main() {
	p := portfolio.NewPortfolio("Main")

	p.Load("./examples/profile.yml")

	p.Refresh()

	fmt.Println("Hello, Portfolio", p.Holdings["VTI"].Status.Value)
}
