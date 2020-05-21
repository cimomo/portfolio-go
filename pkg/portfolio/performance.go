package portfolio

import (
	"math"
	"time"

	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
)

// Performance analyzes the historic performance of a portfolio
type Performance struct {
	Portfolio           *Portfolio
	StartDate           time.Time
	Historic            []Historic
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

// Historic represents a historic quote or portfolio value
type Historic struct {
	Open  float64
	Close float64
	Date  time.Time
}

// NewPerformance creates a new analysis of the historic performance of a portfolio
func NewPerformance(portfolio *Portfolio) *Performance {
	return &Performance{
		Portfolio: portfolio,
	}
}

// Compute generates the performance data for the portfolio
func (performance *Performance) Compute() error {
	startDate, err := computeStartDateForPortfolio(performance.Portfolio)
	if err != nil {
		return err
	}

	performance.StartDate = startDate

	monthly, err := computeMonthlyBalances(performance.Portfolio, performance.StartDate)
	if err != nil {
		return err
	}

	performance.Historic = monthly
	performance.InitialBalance = monthly[0].Open
	performance.FinalBalance = monthly[len(monthly)-1].Close

	cagr, err := computeCAGR(performance.StartDate, performance.InitialBalance, performance.FinalBalance)
	if err != nil {
		return err
	}

	performance.CAGR = cagr

	sd := computeStandardDeviation(performance.Historic)

	performance.Stdev = sd

	return nil
}

func computeStartDateForPortfolio(portfolio *Portfolio) (time.Time, error) {
	now := time.Now()
	thisYear := now.Year()
	startYear := thisYear - 10
	earliestDate := &datetime.Datetime{
		Day:   1,
		Month: 1,
		Year:  startYear,
	}
	earliest := *earliestDate.Time()
	startDate := earliest

	for _, symbol := range portfolio.Symbols {
		start, err := computeStartDateForAsset(earliest, symbol)
		if err != nil {
			return time.Time{}, err
		}

		if (start).After(startDate) {
			startDate = start
		}
	}

	return startDate, nil
}

func computeStartDateForAsset(earliest time.Time, symbol string) (time.Time, error) {
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

func computeCAGR(startDate time.Time, initialBalance float64, finalBalance float64) (float64, error) {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		return 0, err
	}

	now := time.Now().In(ny)

	duration := now.Sub(startDate)
	hours := duration.Hours()
	years := hours / 24 / 365

	cagr := math.Pow(finalBalance/initialBalance, 1/years) - 1

	return cagr, nil
}

func computeMonthlyBalancesForAsset(symbol string, startDate time.Time) ([]finance.ChartBar, error) {
	monthly := make([]finance.ChartBar, 0)

	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		return nil, err
	}

	now := time.Now().In(ny)

	p := &chart.Params{
		Symbol:   symbol,
		Start:    datetime.New(&startDate),
		End:      datetime.New(&now),
		Interval: datetime.OneMonth,
	}

	iter := chart.Get(p)
	if iter.Err() != nil {
		return nil, err
	}

	for iter.Next() {
		b := iter.Bar()
		monthly = append(monthly, *b)
	}

	return monthly, nil
}

func computeMonthlyBalances(portfolio *Portfolio, startDate time.Time) ([]Historic, error) {
	var monthly []Historic

	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		return nil, err
	}

	for _, symbol := range portfolio.Symbols {
		holding := portfolio.Holdings[symbol]
		monthlyForAsset, err := computeMonthlyBalancesForAsset(symbol, startDate)
		if err != nil {
			return nil, err
		}

		if monthly == nil {
			monthly = make([]Historic, len(monthlyForAsset))
		}

		for i := range monthlyForAsset {
			open, _ := monthlyForAsset[i].Open.Float64()
			monthly[i].Open += open * holding.Quantity

			close, _ := monthlyForAsset[i].AdjClose.Float64()
			monthly[i].Close += close * holding.Quantity

			monthly[i].Date = time.Unix(int64(monthlyForAsset[i].Timestamp), 0).In(ny)
		}
	}

	return monthly, nil
}

func computeStandardDeviation(historic []Historic) float64 {
	returns := computeMonthlyReturns(historic)

	var sum, mean, sd float64
	n := float64(len(returns))

	for _, r := range returns {
		sum += r
	}

	mean = sum / n

	for _, r := range returns {
		sd += math.Pow(r-mean, 2)
	}

	// Annualized standard deviation of monthly returns
	sd = math.Sqrt(sd/n) * math.Sqrt(12)

	return sd
}

func computeMonthlyReturns(historic []Historic) []float64 {
	returns := make([]float64, len(historic))

	for i, month := range historic {
		if i == 0 {
			returns[i] = ((month.Close - month.Open) / month.Open) * 100
		} else {
			returns[i] = ((historic[i].Close - historic[i-1].Close) / historic[i-1].Close) * 100
		}
	}

	return returns
}
