package portfolio

import (
	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/index"
)

// Market index symbols
const (
	Dow         string = "^DJI"
	SP500       string = "^GSPC"
	Nasdaq      string = "^IXIC"
	Russell2000 string = "^RUT"
)

// Market defines the broader market indices we track
type Market struct {
	Dow         *finance.Index
	SP500       *finance.Index
	Nasdaq      *finance.Index
	Russell2000 *finance.Index
}

// NewMarket returns a new market
func NewMarket() *Market {
	return &Market{}
}

// Refresh fetches the latest quotes for the market indices
func (market *Market) Refresh() error {
	dow, err := index.Get(Dow)
	if err != nil {
		return err
	}

	sp500, err := index.Get(SP500)
	if err != nil {
		return err
	}

	nasdaq, err := index.Get(Nasdaq)
	if err != nil {
		return err
	}

	russell2000, err := index.Get(Russell2000)
	if err != nil {
		return err
	}

	market.Dow = dow
	market.SP500 = sp500
	market.Nasdaq = nasdaq
	market.Russell2000 = russell2000

	return nil
}
