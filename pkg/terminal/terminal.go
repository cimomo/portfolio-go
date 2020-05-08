package terminal

import (
	"time"

	"github.com/cimomo/portfolio-go/pkg/portfolio"
	"github.com/rivo/tview"
)

// Terminal defines the main terminal window for portfolio visualization
type Terminal struct {
	application       *tview.Application
	market            *portfolio.Market
	portfolio         *portfolio.Portfolio
	performance       *portfolio.Performance
	marketViewer      *MarketViewer
	portfolioViewer   *PortfolioViewer
	performanceViewer *PerformanceViewer
}

// NewTerminal returns a new terminal window
func NewTerminal(market *portfolio.Market, portfolio *portfolio.Portfolio, performance *portfolio.Performance) *Terminal {
	return &Terminal{
		application: tview.NewApplication(),
		portfolio:   portfolio,
		market:      market,
		performance: performance,
	}
}

// Start starts the terminal application
func (term *Terminal) Start() {
	portfolioViewer := NewPortfolioViewer(term.portfolio)
	term.portfolioViewer = portfolioViewer

	marketViewer := NewMarketViewer(term.market)
	term.marketViewer = marketViewer

	performanceViewer := NewPerformanceViewer(term.performance)
	term.performanceViewer = performanceViewer

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
				term.marketViewer.Refresh()
			})
		}
	}
}

func (term *Terminal) setLayout() {
	grid := tview.NewGrid().SetRows(4, 0, 16).SetColumns(0).SetBorders(false).
		AddItem(term.marketViewer.table, 0, 0, 1, 1, 0, 0, false).
		AddItem(term.portfolioViewer.table, 1, 0, 1, 1, 0, 0, false).
		AddItem(term.performanceViewer.table, 2, 0, 1, 1, 0, 0, false)
	term.application.SetRoot(grid, true)
}
