package portfolio

import (
	"math"
	"time"

	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
	"github.com/piquette/finance-go/quote"
)

// Performance analyzes the historic performance of a portfolio and compares it against a benchmark
type Performance struct {
	Portfolio       *Portfolio
	BenchmarkSymbol string
	StartDate       time.Time
	EndDate         time.Time
	Result          *PerformanceResult
	Benchmark       *PerformanceResult
}

// PerformanceResult contains the historic performance of a portfolio
type PerformanceResult struct {
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
func NewPerformance(portfolio *Portfolio, benchmark string) *Performance {
	return &Performance{
		Portfolio:       portfolio,
		BenchmarkSymbol: benchmark,
		Result:          NewPerformanceResult(),
		Benchmark:       NewPerformanceResult(),
	}
}

// NewPerformanceResult creates a new historic performance result of a portfolio
func NewPerformanceResult() *PerformanceResult {
	return &PerformanceResult{}
}

// Compute generates the performance data for the portfolio
func (performance *Performance) Compute() error {
	startDate, endDate, err := computeStartAndEndDateForPortfolio(performance.Portfolio)
	if err != nil {
		return err
	}
	performance.StartDate = startDate
	performance.EndDate = endDate

	monthly, err := computeMonthlyBalances(performance.Portfolio, performance.StartDate, performance.EndDate)
	if err != nil {
		return err
	}
	performance.Result.Historic = monthly
	performance.Result.InitialBalance = monthly[0].Open
	performance.Result.FinalBalance = monthly[len(monthly)-1].Close

	cagr, err := computeCAGR(performance.StartDate, performance.EndDate, performance.Result.InitialBalance, performance.Result.FinalBalance)
	if err != nil {
		return err
	}
	performance.Result.CAGR = cagr

	monthlyReturns := computeMonthlyReturns(performance.Result.Historic)

	sd := computeStandardDeviation(monthlyReturns)
	performance.Result.Stdev = sd

	maxDrawdown := computeMaxDrawdown(monthlyReturns)
	performance.Result.MaxDrawdown = maxDrawdown

	yearly := computeYearlyReturns(performance.Result.Historic, performance.StartDate, performance.EndDate)
	best, worst := computeBestAndWorstYears(yearly)
	performance.Result.BestYear = best
	performance.Result.WorstYear = worst

	sharpe, err := computeSharpeRatio(performance.Result.CAGR, performance.Result.Stdev)
	if err != nil {
		return err
	}
	performance.Result.SharpeRatio = sharpe

	return nil
}

func computeStartAndEndDateForPortfolio(portfolio *Portfolio) (time.Time, time.Time, error) {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	now := time.Now().In(ny)
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
		start, err := computeStartDateForAsset(earliest, now, symbol)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}

		if (start).After(startDate) {
			startDate = start
		}
	}

	return startDate, now, nil
}

func computeStartDateForAsset(earliest time.Time, endDate time.Time, symbol string) (time.Time, error) {
	p := &chart.Params{
		Symbol:   symbol,
		Start:    datetime.New(&earliest),
		End:      datetime.New(&endDate),
		Interval: datetime.OneDay,
	}

	iter := chart.Get(p)
	for iter.Next() {
		b := iter.Bar()
		startDate := time.Unix(int64(b.Timestamp), 0).In(earliest.Location())
		return startDate, nil
	}

	return time.Time{}, iter.Err()
}

func computeCAGR(startDate time.Time, endDate time.Time, initialBalance float64, finalBalance float64) (float64, error) {
	duration := endDate.Sub(startDate)
	hours := duration.Hours()
	years := hours / 24 / 365

	cagr := (math.Pow(finalBalance/initialBalance, 1/years) - 1) * 100

	return cagr, nil
}

func computeMonthlyBalancesForAsset(symbol string, startDate time.Time, endDate time.Time) ([]finance.ChartBar, error) {
	monthly := make([]finance.ChartBar, 0)

	p := &chart.Params{
		Symbol:   symbol,
		Start:    datetime.New(&startDate),
		End:      datetime.New(&endDate),
		Interval: datetime.OneMonth,
	}

	iter := chart.Get(p)
	if iter.Err() != nil {
		return nil, iter.Err()
	}

	for iter.Next() {
		b := iter.Bar()
		monthly = append(monthly, *b)
	}

	return monthly, nil
}

func computeMonthlyBalances(portfolio *Portfolio, startDate time.Time, endDate time.Time) ([]Historic, error) {
	var monthly []Historic

	for _, symbol := range portfolio.Symbols {
		holding := portfolio.Holdings[symbol]
		monthlyForAsset, err := computeMonthlyBalancesForAsset(symbol, startDate, endDate)
		if err != nil {
			return nil, err
		}

		if monthly == nil {
			monthly = make([]Historic, len(monthlyForAsset))
		}

		for i := range monthly {
			open, _ := monthlyForAsset[i].Open.Float64()
			monthly[i].Open += open * holding.Quantity

			close, _ := monthlyForAsset[i].AdjClose.Float64()
			monthly[i].Close += close * holding.Quantity

			monthly[i].Date = time.Unix(int64(monthlyForAsset[i].Timestamp), 0).In(startDate.Location())
		}
	}

	return monthly, nil
}

func computeStandardDeviation(monthlyReturns []float64) float64 {
	var sum, mean, sd float64
	n := float64(len(monthlyReturns))

	for _, r := range monthlyReturns {
		sum += r
	}

	mean = sum / n

	for _, r := range monthlyReturns {
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

func computeYearlyReturns(historic []Historic, startDate time.Time, endDate time.Time) []float64 {
	startYear := startDate.Year()
	endYear := endDate.Year()
	years := endYear - startYear + 1

	returns := make([]float64, years)

	lastClose := historic[0].Open
	curr := 0

	for _, month := range historic {
		currYear := startYear + curr
		if month.Date.Year() == currYear && month.Date.Month() == time.December {
			returns[curr] = ((month.Close - lastClose) / lastClose) * 100
			lastClose = month.Close
			curr++
		}
	}

	returns[years-1] = ((historic[len(historic)-1].Close - lastClose) / lastClose) * 100

	return returns
}

func computeBestAndWorstYears(yearlyReturns []float64) (float64, float64) {
	var best, worst float64
	for _, yearly := range yearlyReturns {
		if yearly > best {
			best = yearly
		} else if yearly < worst {
			worst = yearly
		}
	}

	return best, worst
}

func computeMaxDrawdown(monthlyReturns []float64) float64 {
	var maxDrawdown float64

	for _, monthly := range monthlyReturns {
		if monthly < maxDrawdown {
			maxDrawdown = monthly
		}
	}

	return maxDrawdown
}

func computeRiskFreeReturn() (float64, error) {
	// We use the yield of the 13-week treasury bill as the risk-free return
	quote, err := quote.Get("^IRX")
	if err != nil {
		return 0, err
	}

	return quote.RegularMarketPrice, nil
}

func computeSharpeRatio(cagr float64, stdev float64) (float64, error) {
	riskFree, err := computeRiskFreeReturn()
	if err != nil {
		return 0, err
	}

	sharpe := (cagr - riskFree) / stdev

	return sharpe, nil
}
