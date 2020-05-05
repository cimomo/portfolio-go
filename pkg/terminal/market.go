package terminal

import (
	"math"

	"github.com/cimomo/portfolio-go/pkg/portfolio"
	"github.com/gdamore/tcell"
	"github.com/piquette/finance-go"
	"github.com/rivo/tview"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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

	viewer.Refresh()

	return &viewer
}

// Refresh fetches the latest market data and refreshes the viewer
func (viewer *MarketViewer) Refresh() {
	viewer.market.Refresh()
	viewer.drawMarket()
}

func (viewer *MarketViewer) drawMarket() {
	market := viewer.market

	viewer.drawIndex("Dow 30", market.Dow, 0)
	viewer.drawIndex("S&P 500", market.SP500, 1)
	viewer.drawIndex("Nasdaq", market.Nasdaq, 2)
	viewer.drawIndex("Russell 2000", market.Russell2000, 3)
	viewer.drawIndex("Foreign", market.Foreign, 4)
	viewer.drawIndex("China", market.China, 5)
	viewer.drawIndex("US Bond", market.USBond, 6)
	viewer.drawIndex("10-Yr Yield", market.Treasury10, 7)
	viewer.drawIndex("Gold", market.Gold, 8)
}

func (viewer *MarketViewer) drawIndex(name string, index *finance.Index, c int) {
	value := index.RegularMarketPrice
	change := index.RegularMarketChange
	percent := index.RegularMarketChangePercent

	bg := tcell.ColorDarkGreen
	formatter := " +%.2f (+%.2f%%)"

	if change < 0 {
		bg = tcell.ColorDarkRed
		formatter = " -%.2f (-%.2f%%)"
	}

	cell := tview.NewTableCell(name).SetTextColor(tcell.ColorYellow).SetBackgroundColor(bg).SetAttributes(tcell.AttrBold).SetAlign(tview.AlignCenter)
	viewer.table.SetCell(0, c, cell)

	printer := message.NewPrinter(language.English)
	dayValue := printer.Sprintf("%.2f", value)
	cell = tview.NewTableCell(dayValue).SetTextColor(tcell.ColorWhite).SetBackgroundColor(bg).SetAlign(tview.AlignCenter)
	viewer.table.SetCell(1, c, cell)

	printer = message.NewPrinter(language.English)
	dayChange := printer.Sprintf(formatter, math.Abs(change), math.Abs(percent))
	cell = tview.NewTableCell(dayChange).SetTextColor(tcell.ColorWhite).SetBackgroundColor(bg).SetAlign(tview.AlignCenter)
	viewer.table.SetCell(2, c, cell)
}
