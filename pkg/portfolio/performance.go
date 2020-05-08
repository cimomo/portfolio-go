package portfolio

import "time"

// Performance analyzes the historic performance of a portfolio
type Performance struct {
	Portfolio           *Portfolio
	StartDate           *time.Time
	InitialBalance      float64
	FinalBalance        float64
	CAGR                float64
	Stdev               float64
	BestYear            float64
	WorstYear           float64
	MaxDrawdown         float64
	SharpeRatio         float64
	USMarketCorrelation float64
}

// NewPerformance creates a new analysis of the historic performance of a portfolio
func NewPerformance(portfolio *Portfolio) *Performance {
	return &Performance{
		Portfolio: portfolio,
	}
}
