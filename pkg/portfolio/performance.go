package portfolio

import (
	"time"

	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
)

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

// Compute generates the performance data for the portfolio
func (performance *Performance) Compute() error {
	startDate, err := performance.getStartDateForPortfolio()
	if err != nil {
		return err
	}

	performance.StartDate = startDate

	return nil
}

func (performance *Performance) getStartDateForPortfolio() (*time.Time, error) {
	now := time.Now()
	thisYear := now.Year()
	startYear := thisYear - 9
	earliestDate := &datetime.Datetime{
		Day:   1,
		Month: 1,
		Year:  startYear,
	}
	earliest := earliestDate.Time()
	startDate := earliest

	for _, symbol := range performance.Portfolio.Symbols {
		start, err := performance.getStartDateForAsset(earliest, symbol)
		if err != nil {
			return nil, err
		}

		if (start).After(*startDate) {
			startDate = start
		}
	}

	return startDate, nil
}

func (performance *Performance) getStartDateForAsset(earliest *time.Time, symbol string) (*time.Time, error) {
	now := time.Now()
	p := &chart.Params{
		Symbol:   symbol,
		Start:    datetime.New(earliest),
		End:      datetime.New(&now),
		Interval: datetime.OneMonth,
	}

	iter := chart.Get(p)
	for iter.Next() {
		b := iter.Bar()
		startDate := time.Unix(int64(b.Timestamp), 0)
		return &startDate, nil
	}

	return nil, iter.Err()
}
