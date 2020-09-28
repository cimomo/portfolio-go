package terminal

import (
	"github.com/cimomo/portfolio-go/pkg/portfolio"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// PortfolioViewer displays real-time portfolio data
type PortfolioViewer struct {
	portfolio *portfolio.Portfolio
	table     *tview.Table
}

// NewPortfolioViewer returns a new viewer for the real-time portfolio data
func NewPortfolioViewer(portfolio *portfolio.Portfolio) *PortfolioViewer {
	return &PortfolioViewer{
		portfolio: portfolio,
		table:     tview.NewTable().SetBorders(false),
	}
}

// Reload updates the portfolio data object
func (viewer *PortfolioViewer) Reload(portfolio *portfolio.Portfolio) {
	viewer.portfolio = portfolio
}

// Draw fetches the latest portfolio data and refreshes the viewer
func (viewer *PortfolioViewer) Draw() {
	viewer.table.Clear()
	viewer.drawHeader()
	viewer.drawPortfolio()
}

func (viewer *PortfolioViewer) drawHeader() {
	var cell *tview.TableCell
	header := []string{
		"SYMBOL", "CLASS", "QUANTITY", "PRICE", "WATCH",
		"1-DAY CHANGE$", "1-DAY CHANGE%",
		"VALUE", "1-Day VALUE CHANGE$",
		"UNREALIZED$", "UNREALIZED%",
		"ALLOCATION", "TARGET",
	}

	for c := 0; c < len(header); c++ {
		cell = tview.NewTableCell(header[c]).SetTextColor(tcell.ColorYellow).SetBackgroundColor(tcell.ColorDarkSlateGray).SetAttributes(tcell.AttrBold)
		if c < 2 {
			cell.SetAlign(tview.AlignLeft)
		} else {
			cell.SetAlign((tview.AlignRight))
		}
		viewer.table.SetCell(0, c, cell)
	}
}

func (viewer *PortfolioViewer) drawPortfolio() {
	port := viewer.portfolio
	holdings := port.Holdings

	r := 1
	for _, symbol := range port.Symbols {
		holding := holdings[symbol]
		setString(viewer.table, symbol, r, 0, tcell.ColorWhite, tview.AlignLeft)
		setString(viewer.table, string(holding.Asset.Subclass), r, 1, tcell.ColorWhite, tview.AlignLeft)
		setQuantity(viewer.table, holding.Quantity, r, 2, tview.AlignCenter)
		setDollarAmountAgainstWatch(viewer.table, holding.Quote.RegularMarketPrice, holding.Watch, r, 3)
		setDollarAmount(viewer.table, holding.Watch, r, 4, tcell.ColorWhite)
		setDollarChange(viewer.table, holding.Quote.RegularMarketChange, r, 5)
		setPercentChange(viewer.table, holding.Quote.RegularMarketChangePercent, r, 6)
		setDollarAmount(viewer.table, holding.Status.Value, r, 7, tcell.ColorWhite)
		setDollarChange(viewer.table, holding.Quote.RegularMarketChange*holding.Quantity, r, 8)
		setDollarChange(viewer.table, holding.Status.Unrealized, r, 9)
		setPercentChange(viewer.table, holding.Status.UnrealizedPercent, r, 10)
		setPercent(viewer.table, port.Status.Allocation[symbol], r, 11, tcell.ColorWhite)
		setPercent(viewer.table, port.TargetAllocation[symbol], r, 12, tcell.ColorWhite)

		r++
	}

	setString(viewer.table, "TOTAL", r, 0, tcell.ColorYellow, tview.AlignLeft)
	setPercentChange(viewer.table, port.Status.RegularMarketChangePercent, r, 6)
	setDollarAmount(viewer.table, port.Status.Value, r, 7, tcell.ColorYellow)
	setDollarChange(viewer.table, port.Status.RegularMarketChange, r, 8)
	setDollarChange(viewer.table, port.Status.Unrealized, r, 9)
	setPercentChange(viewer.table, port.Status.UnrealizedPercent, r, 10)
	setPercent(viewer.table, 100.0, r, 11, tcell.ColorYellow)
	setPercent(viewer.table, 100.0, r, 12, tcell.ColorYellow)
}
