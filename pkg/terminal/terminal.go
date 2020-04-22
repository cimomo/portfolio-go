package terminal

import (
	"time"

	"github.com/cimomo/portfolio-go/pkg/portfolio"
	"github.com/rivo/tview"
)

// Terminal defines the main terminal window for portfolio visualization
type Terminal struct {
	application     *tview.Application
	market          *portfolio.Market
	portfolio       *portfolio.Portfolio
	marketViewer    *MarketViewer
	portfolioViewer *PortfolioViewer
}

// NewTerminal returns a new terminal window
func NewTerminal(market *portfolio.Market, portfolio *portfolio.Portfolio) *Terminal {
	return &Terminal{
		application: tview.NewApplication(),
		portfolio:   portfolio,
		market:      market,
	}
}

// Start starts the terminal application
func (term *Terminal) Start() {
	portfolioViewer := NewPortfolioViewer(term.portfolio)
	term.portfolioViewer = portfolioViewer

	marketViewer := NewMarketViewer(term.portfolio)
	term.marketViewer = marketViewer

	term.setLayout()

	go term.refresh()
	term.application.Run()
}

// Stop stops the terminal application
func (term *Terminal) Stop() {
	term.application.Stop()
}

func (term *Terminal) refresh() {
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ticker.C:
			term.application.QueueUpdateDraw(func() {
				term.portfolioViewer.Refresh()
			})
		}
	}
}

func (term *Terminal) setLayout() {
	term.application.SetRoot(term.portfolioViewer.table, true)
}
