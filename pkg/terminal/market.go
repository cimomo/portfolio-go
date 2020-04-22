package terminal

import (
	"github.com/cimomo/portfolio-go/pkg/portfolio"
	"github.com/rivo/tview"
)

// MarketViewer displays real-time market data
type MarketViewer struct {
	portfolio *portfolio.Portfolio
	table     *tview.Table
}

// NewMarketViewer returns a new viewer for the real-time market data
func NewMarketViewer(portfolio *portfolio.Portfolio) *MarketViewer {
	viewer := MarketViewer{
		portfolio: portfolio,
		table:     tview.NewTable().SetBorders(false),
	}

	return &viewer
}
