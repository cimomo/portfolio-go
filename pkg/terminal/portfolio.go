package terminal

import (
	"math"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// PortfolioViewer displays real-time portfolio data
type PortfolioViewer struct {
	terminal *Terminal
	table    *tview.Table
}

// NewPortfolioViewer returns a new viewer for the real-time portfolio data
func NewPortfolioViewer(term *Terminal) *PortfolioViewer {
	viewer := PortfolioViewer{
		terminal: term,
		table:    tview.NewTable().SetBorders(false),
	}

	viewer.drawHeader()
	viewer.drawPortfolio()

	return &viewer
}

// Connect adds the portfolio viewer to the parent terminal
func (viewer *PortfolioViewer) Connect() {
	viewer.terminal.application.SetRoot(viewer.table, true)
}

// Refresh fetches the latest portfolio data and refreshes the viewer
func (viewer *PortfolioViewer) Refresh() {
	viewer.terminal.portfolio.Refresh()
	viewer.drawPortfolio()
}

func (viewer *PortfolioViewer) drawHeader() {
	var cell *tview.TableCell
	header := []string{
		"SYMBOL", "CLASS", "QUANTITY", "PRICE",
		"1-Day CHANGE$", "1-Day CHANGE%",
		"VALUE", "1-Day VALUE CHANGE$",
		"UNREALIZED GAIN/LOSS$", "UNREALIZED GAIN/LOSS%",
		"Allocation", "Target",
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
	port := viewer.terminal.portfolio
	holdings := port.Holdings

	r := 1
	for symbol, holding := range holdings {
		viewer.setString(symbol, r, 0, tcell.ColorWhite, tview.AlignLeft)
		viewer.setString(string(holding.Asset.Subclass), r, 1, tcell.ColorWhite, tview.AlignLeft)
		viewer.setQuantity(holding.Quantity, r, 2)
		viewer.setDollarAmount(holding.Quote.RegularMarketPrice, r, 3, tcell.ColorWhite)
		viewer.setDollarChange(holding.Quote.RegularMarketChange, r, 4)
		viewer.setPercentChange(holding.Quote.RegularMarketChangePercent, r, 5)
		viewer.setDollarAmount(holding.Status.Value, r, 6, tcell.ColorWhite)
		viewer.setDollarChange(holding.Quote.RegularMarketChange*holding.Quantity, r, 7)
		viewer.setDollarChange(holding.Status.Unrealized, r, 8)
		viewer.setPercentChange(holding.Status.UnrealizedPercent, r, 9)
		viewer.setPercent(port.Status.Allocation[symbol], r, 10, tcell.ColorWhite)
		viewer.setPercent(port.TargetAllocation[symbol], r, 11, tcell.ColorWhite)

		r++
	}

	viewer.setString("TOTAL", r, 0, tcell.ColorYellow, tview.AlignLeft)
	viewer.setPercentChange(port.Status.RegularMarketChangePercent, r, 5)
	viewer.setDollarAmount(port.Status.Value, r, 6, tcell.ColorYellow)
	viewer.setDollarChange(port.Status.RegularMarketChange, r, 7)
	viewer.setDollarChange(port.Status.Unrealized, r, 8)
	viewer.setPercentChange(port.Status.UnrealizedPercent, r, 9)
	viewer.setPercent(100.0, r, 10, tcell.ColorYellow)
	viewer.setPercent(100.0, r, 11, tcell.ColorYellow)
}

func (viewer *PortfolioViewer) setPercentChange(value float64, r int, c int) {
	color := tcell.ColorGreen
	if value < 0 {
		color = tcell.ColorRed
	}
	viewer.setPercent(value, r, c, color)
}

func (viewer *PortfolioViewer) setDollarChange(value float64, r int, c int) {
	color := tcell.ColorGreen
	if value < 0 {
		color = tcell.ColorRed
	}
	viewer.setDollarAmount(value, r, c, color)
}

func (viewer *PortfolioViewer) setPercent(value float64, r int, c int, color tcell.Color) {
	viewer.setFloat64(value, "%.2f%%", r, c, color)
}

func (viewer *PortfolioViewer) setDollarAmount(value float64, r int, c int, color tcell.Color) {
	formatter := "$%.2f"
	if value < 0 {
		formatter = "-$%.2f"
	}

	v := math.Abs(value)
	viewer.setFloat64(v, formatter, r, c, color)
}

func (viewer *PortfolioViewer) setQuantity(value float64, r int, c int) {
	viewer.setFloat64(value, "%.2f", r, c, tcell.ColorWhite)
}

func (viewer *PortfolioViewer) setFloat64(value float64, formatter string, r int, c int, color tcell.Color) {
	printer := message.NewPrinter(language.English)
	fValue := printer.Sprintf(formatter, value)
	viewer.setString(fValue, r, c, color, tview.AlignRight)
}

func (viewer *PortfolioViewer) setString(value string, r int, c int, color tcell.Color, align int) {
	cell := tview.NewTableCell(value).SetTextColor(color).SetAlign(align).SetExpansion(1)
	viewer.table.SetCell(r, c, cell)
}
