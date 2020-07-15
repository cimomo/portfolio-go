package portfolio

import (
	"errors"

	"github.com/piquette/finance-go/quote"
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

type portfolioConfig struct {
	Name     string          `yaml:"portfolio"`
	Holdings []holdingConfig `yaml:"holdings"`
}

// NewPortfolio returns an empty portfolio of asset holdings
func NewPortfolio() *Portfolio {
	return &Portfolio{
		Symbols:          make([]string, 0),
		Holdings:         make(map[string]*Holding),
		TargetAllocation: make(map[string]float64),
		Status:           &Status{},
	}
}

// Load loads a portfolio from the config
func (portfolio *Portfolio) Load(config portfolioConfig) error {

	portfolio.Name = config.Name

	totalAllocation := 0.0

	for _, holdingConfig := range config.Holdings {
		portfolio.Symbols = append(portfolio.Symbols, holdingConfig.Symbol)
		portfolio.Holdings[holdingConfig.Symbol] = NewHolding(
			holdingConfig.Symbol,
			holdingConfig.Quantity,
			holdingConfig.CostBasis)
		portfolio.TargetAllocation[holdingConfig.Symbol] = holdingConfig.TargetAllocation
		totalAllocation += holdingConfig.TargetAllocation
		portfolio.CostBasis += holdingConfig.CostBasis
	}

	if totalAllocation != 100.0 && totalAllocation != 0.0 {
		return errors.New("Total allocation should be either 0% (ignored) or 100%")
	}

	return nil
}

// Refresh computes the current status of the entire portfolio and its holdings
func (portfolio *Portfolio) Refresh() error {
	result := quote.List(portfolio.Symbols)

	if result.Err() != nil {
		return result.Err()
	}

	for result.Next() {
		quote := result.Quote()
		symbol := quote.Symbol
		holding := portfolio.Holdings[symbol]
		holding.Quote = quote
		holding.RefreshStatus()
	}

	portfolio.RefreshStatus()

	return nil
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

	status.RegularMarketChangePercent = 0
	if previousValue != 0 {
		status.RegularMarketChangePercent = (status.RegularMarketChange / previousValue) * 100
	}

	status.UnrealizedPercent = 0
	if portfolio.CostBasis != 0 {
		status.UnrealizedPercent = (status.Unrealized / portfolio.CostBasis) * 100
	}

	for symbol, holding := range portfolio.Holdings {
		status.Allocation[symbol] = 0
		if status.Value != 0 {
			status.Allocation[symbol] = (holding.Status.Value / status.Value) * 100
		}
	}

	portfolio.Status = &status
}

// Clone makes a copy of the Portfolio
func (portfolio *Portfolio) Clone() *Portfolio {
	port := Portfolio{
		Name:      portfolio.Name,
		CostBasis: portfolio.CostBasis,
		Status:    &Status{},
	}

	symbols := make([]string, len(portfolio.Symbols))
	for i := range portfolio.Symbols {
		symbols[i] = portfolio.Symbols[i]
	}

	holdings := make(map[string]*Holding)
	for k, v := range portfolio.Holdings {
		holdings[k] = v.Clone()
	}

	allocation := make(map[string]float64)
	for k, v := range portfolio.TargetAllocation {
		allocation[k] = v
	}

	port.Symbols = symbols
	port.Holdings = holdings
	port.TargetAllocation = allocation

	return &port
}
