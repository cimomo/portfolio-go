package portfolio

import (
	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/quote"
)

// Holding defines a position in the portfolio
type Holding struct {
	Asset            *Asset
	Quote            *finance.Quote
	TargetAllocation float64
	Quantity         float64
	CostBasis        float64
	Status           *HoldingStatus
}

// HoldingStatus defines the real-time status of a particular holding
type HoldingStatus struct {
	Value             float64
	Unrealized        float64
	UnrealizedPercent float64
}

// NewHolding returns a new holding object
func NewHolding(symbol string, allocation float64, quantity float64, basis float64) *Holding {
	holding := Holding{}

	asset := NewAsset(symbol)

	if asset == nil {
		return nil
	}

	holding.Asset = asset
	holding.TargetAllocation = allocation
	holding.Quantity = quantity
	holding.CostBasis = basis

	return &holding
}

// Refresh gets the current quote and computes the current status of a particular holding
func (holding *Holding) Refresh() error {
	err := holding.RefreshQuote()
	if err != nil {
		return err
	}

	holding.RefreshStatus()

	return nil
}

// RefreshQuote gets the current quote of a particular holding
func (holding *Holding) RefreshQuote() error {
	quote, err := quote.Get(holding.Asset.Symbol)
	if err != nil {
		return err
	}

	holding.Quote = quote

	return nil
}

// RefreshStatus computes the current status of a holding from the current quote
func (holding *Holding) RefreshStatus() {
	quote := holding.Quote
	status := HoldingStatus{}

	status.Value = quote.RegularMarketPrice * holding.Quantity
	status.Unrealized = status.Value - holding.CostBasis
	status.UnrealizedPercent = (status.Unrealized / holding.CostBasis) * 100

	holding.Status = &status
}
