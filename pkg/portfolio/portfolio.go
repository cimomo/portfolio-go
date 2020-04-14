package portfolio

import (
	"io/ioutil"

	"github.com/piquette/finance-go/quote"
	"gopkg.in/yaml.v2"
)

// Portfolio defines a portfolio of asset holdings
type Portfolio struct {
	Name     string
	Holdings map[string]*Holding
}

type holdingConfig struct {
	Symbol           string  `yaml:"symbol"`
	TargetAllocation float64 `yaml:"allocation"`
	Quantity         float64 `yaml:"quantity"`
	CostBasis        float64 `yaml:"basis"`
}

type portfolioConfig []holdingConfig

// NewPortfolio returns an empty portfolio of asset holdings
func NewPortfolio(name string) *Portfolio {
	return &Portfolio{
		Name:     name,
		Holdings: make(map[string]*Holding),
	}
}

// Load loads a portfolio from the given file
func (portfolio *Portfolio) Load(profile string) error {
	file, err := ioutil.ReadFile(profile)
	if err != nil {
		return err
	}

	portfolioConfig := portfolioConfig{}

	err = yaml.Unmarshal(file, &portfolioConfig)
	if err != nil {
		return err
	}

	for _, holdingConfig := range portfolioConfig {
		portfolio.Holdings[holdingConfig.Symbol] = NewHolding(
			holdingConfig.Symbol,
			holdingConfig.TargetAllocation,
			holdingConfig.Quantity,
			holdingConfig.CostBasis)
	}

	return nil
}

// Refresh computes the current status of the entire portfolio
func (portfolio *Portfolio) Refresh() {
	symbols := make([]string, 0, len(portfolio.Holdings))

	for symbol := range portfolio.Holdings {
		symbols = append(symbols, symbol)
	}

	result := quote.List(symbols)

	for result.Next() {
		quote := result.Quote()
		symbol := quote.Symbol
		holding := portfolio.Holdings[symbol]
		holding.Quote = quote
		holding.RefreshStatus()
	}
}
