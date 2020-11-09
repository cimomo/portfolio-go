package portfolio

import (
	"fmt"

	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/index"
)

// Market index symbols
const (
	Dow         string = "^DJI"
	SP500       string = "^GSPC"
	Nasdaq      string = "^IXIC"
	Russell2000 string = "^RUT"
	Foreign     string = "VXUS"
	China       string = "000001.SS"
	USBond      string = "BND"
	Treasury10  string = "^TNX"
	Gold        string = "GC=F"
	Oil         string = "CL=F"
	Bitcoin     string = "BTCUSD=X"
	Ethereum    string = "ETHUSD=X"
)

// Market defines the broader market indices we track
type Market struct {
	Dow         *finance.Index
	SP500       *finance.Index
	Nasdaq      *finance.Index
	Russell2000 *finance.Index
	Foreign     *finance.Index
	China       *finance.Index
	USBond      *finance.Index
	Treasury10  *finance.Index
	Gold        *finance.Index
	Oil         *finance.Index
	Bitcoin     *finance.Index
	Ethereum    *finance.Index
}

// NewMarket returns a new market
func NewMarket() *Market {
	return &Market{}
}

// Refresh fetches the latest quotes for the market indices
func (market *Market) Refresh() error {
	indices := []string{
		Dow, SP500, Nasdaq, Russell2000, Foreign, China, USBond, Treasury10, Gold, Oil, Bitcoin, Ethereum,
	}

	result := index.List(indices)

	if result.Err() != nil {
		return result.Err()
	}

	for result.Next() {
		index := result.Index()

		switch index.Symbol {
		case Dow:
			market.Dow = index
		case SP500:
			market.SP500 = index
		case Nasdaq:
			market.Nasdaq = index
		case Russell2000:
			market.Russell2000 = index
		case Foreign:
			market.Foreign = index
		case China:
			market.China = index
		case USBond:
			market.USBond = index
		case Treasury10:
			market.Treasury10 = index
		case Gold:
			market.Gold = index
		case Oil:
			market.Oil = index
		case Bitcoin:
			market.Bitcoin = index
		case Ethereum:
			market.Ethereum = index
		default:
			return fmt.Errorf("Unknown symbol: %s", index.Symbol)
		}
	}

	return nil
}
