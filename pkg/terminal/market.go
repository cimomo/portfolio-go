package terminal

import (
	"github.com/cimomo/portfolio-go/pkg/portfolio"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// MarketViewer displays real-time market data
type MarketViewer struct {
	market *portfolio.Market
	table  *tview.Table
}

// NewMarketViewer returns a new viewer for the real-time market data
func NewMarketViewer(market *portfolio.Market) *MarketViewer {
	viewer := MarketViewer{
		market: market,
		table:  tview.NewTable().SetBorders(false).SetSeparator(' '),
	}

	viewer.drawHeader()
	viewer.Refresh()

	return &viewer
}

// Refresh fetches the latest market data and refreshes the viewer
func (viewer *MarketViewer) Refresh() {
	viewer.market.Refresh()
	viewer.drawMarket()
}

func (viewer *MarketViewer) drawHeader() {
	var cell *tview.TableCell
	header := []string{
		"Dow", "S&P 500", "Nasdaq",
	}

	for c := 0; c < len(header); c++ {
		cell = tview.NewTableCell(header[c]).
			SetTextColor(tcell.ColorYellow).
			SetBackgroundColor(tcell.ColorDarkSlateGray).
			SetAttributes(tcell.AttrBold).
			SetAlign(tview.AlignCenter)

		viewer.table.SetCell(0, c, cell)
	}
}

func (viewer *MarketViewer) drawMarket() {
	market := viewer.market

	setQuantity(viewer.table, market.Dow.RegularMarketPrice, 1, 0, tview.AlignCenter, 0)
	setQuantity(viewer.table, market.SP500.RegularMarketPrice, 1, 1, tview.AlignCenter, 0)
	setQuantity(viewer.table, market.Nasdaq.RegularMarketPrice, 1, 2, tview.AlignCenter, 0)

	setChangeAndPercent(viewer.table, market.Dow.RegularMarketChange, market.Dow.RegularMarketChangePercent, 2, 0, tview.AlignCenter, 0)
	setChangeAndPercent(viewer.table, market.SP500.RegularMarketChange, market.Dow.RegularMarketChangePercent, 2, 1, tview.AlignCenter, 0)
	setChangeAndPercent(viewer.table, market.Nasdaq.RegularMarketChange, market.Dow.RegularMarketChangePercent, 2, 2, tview.AlignCenter, 0)
}
