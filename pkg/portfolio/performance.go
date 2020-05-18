package portfolio

import (
	"errors"
	"math"
	"time"

	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
)

// Performance analyzes the historic performance of a portfolio
type Performance struct {
	Portfolio           *Portfolio
	StartDate           time.Time
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
	startDate, err := performance.computeStartDateForPortfolio()
	if err != nil {
		return err
	}

	performance.StartDate = startDate

	initialBalance, err := performance.computeInitialBalance(startDate)
	if err != nil {
		return err
	}

	performance.InitialBalance = initialBalance

	finalBalance, err := performance.computeFinalBalance()
	if err != nil {
		return err
	}

	performance.FinalBalance = finalBalance

	cagr, err := performance.computeCAGR()
	if err != nil {
		return err
	}

	performance.CAGR = cagr

	return nil
}

func (performance *Performance) computeStartDateForPortfolio() (time.Time, error) {
	now := time.Now()
	thisYear := now.Year()
	startYear := thisYear - 9
	earliestDate := &datetime.Datetime{
		Day:   1,
		Month: 1,
		Year:  startYear,
	}
	earliest := *earliestDate.Time()
	startDate := earliest

	for _, symbol := range performance.Portfolio.Symbols {
		start, err := performance.computeStartDateForAsset(earliest, symbol)
		if err != nil {
			return time.Time{}, err
		}

		if (start).After(startDate) {
			startDate = start
		}
	}

	return startDate, nil
}

func (performance *Performance) computeStartDateForAsset(earliest time.Time, symbol string) (time.Time, error) {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		return time.Time{}, err
	}

	now := time.Now().In(ny)
	p := &chart.Params{
		Symbol:   symbol,
		Start:    datetime.New(&earliest),
		End:      datetime.New(&now),
		Interval: datetime.OneDay,
	}

	iter := chart.Get(p)
	for iter.Next() {
		b := iter.Bar()
		startDate := time.Unix(int64(b.Timestamp), 0).In(ny)
		return startDate, nil
	}

	return time.Time{}, iter.Err()
}

func (performance *Performance) computeInitialBalance(startDate time.Time) (float64, error) {
	var initialBalance float64

	for _, holding := range performance.Portfolio.Holdings {
		bar, err := holding.GetHistoricalQuote(startDate)
		if err != nil {
			return 0, err
		}

		adjClose, _ := bar.AdjClose.Float64()
		initialBalance += adjClose * holding.Quantity
	}

	return initialBalance, nil
}

func (performance *Performance) computeFinalBalance() (float64, error) {
	if performance.Portfolio == nil || performance.Portfolio.Status == nil {
		return 0, errors.New("Portfolio not initialized")
	}

	return performance.Portfolio.Status.Value, nil
}

func (performance *Performance) computeCAGR() (float64, error) {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		return 0, err
	}

	now := time.Now().In(ny)
	start := performance.StartDate

	duration := now.Sub(start)
	hours := duration.Hours()
	years := hours / 24 / 365

	cagr := math.Pow(performance.FinalBalance/performance.InitialBalance, 1/years) - 1

	return cagr, nil
}
