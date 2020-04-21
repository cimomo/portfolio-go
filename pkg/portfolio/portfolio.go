package portfolio

import (
	"errors"
	"io/ioutil"

	"github.com/piquette/finance-go/quote"
	"gopkg.in/yaml.v2"
)

// Portfolio defines a portfolio of asset holdings
type Portfolio struct {
	Name             string
	CostBasis        float64
	Symbols          []string
	Holdings         map[string]*Holding
	TargetAllocation map[string]float64
	Status           *Status
}

// Status defines the real-time status of the entire portfolio
type Status struct {
	Value                      float64
	RegularMarketChange        float64
	RegularMarketChangePercent float64
	Unrealized                 float64
	UnrealizedPercent          float64
	Allocation                 map[string]float64
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
		Name:             name,
		Symbols:          make([]string, 0),
		Holdings:         make(map[string]*Holding),
		TargetAllocation: make(map[string]float64),
		Status:           &Status{},
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

	totalAllocation := 0.0

	for _, holdingConfig := range portfolioConfig {
		portfolio.Symbols = append(portfolio.Symbols, holdingConfig.Symbol)
		portfolio.Holdings[holdingConfig.Symbol] = NewHolding(
			holdingConfig.Symbol,
			holdingConfig.Quantity,
			holdingConfig.CostBasis)
		portfolio.TargetAllocation[holdingConfig.Symbol] = holdingConfig.TargetAllocation
		totalAllocation += holdingConfig.TargetAllocation
		portfolio.CostBasis += holdingConfig.CostBasis
	}

	if totalAllocation != 100.0 {
		return errors.New("Total allocation should be 100%")
	}

	return nil
}

// Refresh computes the current status of the entire portfolio and its holdings
func (portfolio *Portfolio) Refresh() {
	result := quote.List(portfolio.Symbols)

	for result.Next() {
		quote := result.Quote()
		symbol := quote.Symbol
		holding := portfolio.Holdings[symbol]
		holding.Quote = quote
		holding.RefreshStatus()
	}

	portfolio.RefreshStatus()
}

// RefreshStatus computes the current status of the entire portfolio
func (portfolio *Portfolio) RefreshStatus() {
	status := Status{
		Allocation: make(map[string]float64),
	}

	for _, holding := range portfolio.Holdings {
		status.Value += holding.Status.Value
		status.RegularMarketChange += holding.Quote.RegularMarketChange * holding.Quantity
		status.Unrealized += holding.Status.Unrealized
	}

	previousValue := status.Value - status.RegularMarketChange
	status.RegularMarketChangePercent = (status.RegularMarketChange / previousValue) * 100
	status.UnrealizedPercent = (status.Unrealized / portfolio.CostBasis) * 100

	for symbol, holding := range portfolio.Holdings {
		status.Allocation[symbol] = (holding.Status.Value / status.Value) * 100
	}

	portfolio.Status = &status
}
