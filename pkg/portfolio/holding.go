package portfolio

import (
	"time"

	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
	"github.com/piquette/finance-go/quote"
)

// Holding defines a position in the portfolio
type Holding struct {
	Asset     *Asset
	Quantity  float64
	CostBasis float64
	Quote     *finance.Quote
	Status    *HoldingStatus
}

// HoldingStatus defines the real-time status of a particular holding
type HoldingStatus struct {
	Value             float64
	Unrealized        float64
	UnrealizedPercent float64
}

// NewHolding returns a new holding object
func NewHolding(symbol string, quantity float64, basis float64) *Holding {
	return &Holding{
		Asset:     NewAsset(symbol),
		Quantity:  quantity,
		CostBasis: basis,
		Quote:     &finance.Quote{},
		Status:    &HoldingStatus{},
	}
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
	status.UnrealizedPercent = 0
	if holding.CostBasis != 0 {
		status.UnrealizedPercent = (status.Unrealized / holding.CostBasis) * 100
	}

	holding.Status = &status
}

// GetHistoricalQuote returns the quote on a specific day
func (holding *Holding) GetHistoricalQuote(date time.Time) (*finance.ChartBar, error) {
	p := &chart.Params{
		Symbol:   holding.Asset.Symbol,
		Start:    datetime.New(&date),
		End:      datetime.New(&date),
		Interval: datetime.OneDay,
	}

	iter := chart.Get(p)

	for iter.Next() {
		b := iter.Bar()
		return b, nil
	}

	return nil, iter.Err()
}
