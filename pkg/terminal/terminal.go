package terminal

import (
	"fmt"

	"github.com/cimomo/portfolio-go/pkg/portfolio"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Terminal defines the main terminal window for portfolio visualization
type Terminal struct {
	application *tview.Application
	portfolio   *portfolio.Portfolio
}

// NewTerminal returns a new terminal window
func NewTerminal(portfolio *portfolio.Portfolio) *Terminal {
	return &Terminal{
		application: tview.NewApplication(),
		portfolio:   portfolio,
	}
}

// Start starts the terminal application
func (term *Terminal) Start() {
	term.application.Run()
}

// Stop stops the terminal application
func (term *Terminal) Stop() {
	term.application.Stop()
}

// Initialize sets up the terminal screen
func (term *Terminal) Initialize() {
	table := tview.NewTable().SetBorders(false)
	header := []string{
		"SYMBOL", "CLASS", "QUANTITY", "PRICE", "1-Day CHANGE$", "1-Day CHANGE%", "VALUE", "1-Day CHANGE$", "UNREALIZED GAIN/LOSS$", "UNREALIZED GAIN/LOSS%",
	}

	for c := 0; c < len(header); c++ {
		cell := tview.NewTableCell(header[c]).SetTextColor(tcell.ColorYellow)
		if c < 2 {
			cell.SetAlign(tview.AlignLeft)
		} else {
			cell.SetAlign((tview.AlignRight))
		}
		table.SetCell(0, c, cell)
	}

	holdings := term.portfolio.Holdings
	r := 1
	for symbol, holding := range holdings {
		cell := tview.NewTableCell(symbol).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignLeft).SetExpansion(1)
		table.SetCell(r, 0, cell)

		class := fmt.Sprintf("%s", holding.Asset.Subclass)
		cell = tview.NewTableCell(class).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignLeft).SetExpansion(1)
		table.SetCell(r, 1, cell)

		quantity := fmt.Sprintf("%.2f", holding.Quantity)
		cell = tview.NewTableCell(quantity).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignRight).SetExpansion(1)
		table.SetCell(r, 2, cell)

		price := fmt.Sprintf("$%.2f", holding.Quote.RegularMarketPrice)
		cell = tview.NewTableCell(price).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignRight).SetExpansion(1)
		table.SetCell(r, 3, cell)

		color := tcell.ColorGreen
		if holding.Quote.RegularMarketChange < 0 {
			color = tcell.ColorRed
		}

		change := fmt.Sprintf("$%.2f", holding.Quote.RegularMarketChange)
		cell = tview.NewTableCell(change).SetTextColor(color).SetAlign(tview.AlignRight).SetExpansion(1)
		table.SetCell(r, 4, cell)

		changeP := fmt.Sprintf("%.2f%%", holding.Quote.RegularMarketChangePercent)
		cell = tview.NewTableCell(changeP).SetTextColor(color).SetAlign(tview.AlignRight).SetExpansion(1)
		table.SetCell(r, 5, cell)

		value := fmt.Sprintf("$%.2f", holding.Status.Value)
		cell = tview.NewTableCell(value).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignRight).SetExpansion(1)
		table.SetCell(r, 6, cell)

		valueChange := holding.Quote.RegularMarketChange * holding.Quantity
		change = fmt.Sprintf("$%.2f", valueChange)
		cell = tview.NewTableCell(change).SetTextColor(color).SetAlign(tview.AlignRight).SetExpansion(1)
		table.SetCell(r, 7, cell)

		color = tcell.ColorGreen
		if holding.Status.Unrealized < 0 {
			color = tcell.ColorRed
		}

		gain := fmt.Sprintf("$%.2f", holding.Status.Unrealized)
		cell = tview.NewTableCell(gain).SetTextColor(color).SetAlign(tview.AlignRight).SetExpansion(1)
		table.SetCell(r, 8, cell)

		gainP := fmt.Sprintf("%.2f%%", holding.Status.UnrealizedPercent)
		cell = tview.NewTableCell(gainP).SetTextColor(color).SetAlign(tview.AlignRight).SetExpansion(1)
		table.SetCell(r, 9, cell)

		r++
	}

	term.application.SetRoot(table, true)
}
