package terminal

import (
	"time"

	"github.com/cimomo/portfolio-go/pkg/portfolio"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const (
	benchmark      = "SPY"
	initialBalance = 100000.00
)

// Terminal defines the main terminal window for portfolio visualization
type Terminal struct {
	application       *tview.Application
	market            *portfolio.Market
	profile           *portfolio.Profile
	performance       *portfolio.Performance
	marketViewer      *MarketViewer
	portfolioViewer   *PortfolioViewer
	performanceViewer *PerformanceViewer
	returnViewer      *ReturnViewer
	signalPerformance chan int
}

// NewTerminal returns a new terminal window
func NewTerminal(profile *portfolio.Profile) *Terminal {
	market := portfolio.NewMarket()
	performance := portfolio.NewPerformance(profile.Portfolios[0], benchmark, initialBalance)

	return &Terminal{
		application:       tview.NewApplication(),
		profile:           profile,
		market:            market,
		performance:       performance,
		signalPerformance: make(chan int),
	}
}

// Start starts the terminal application
func (term *Terminal) Start() error {
	portfolioViewer := NewPortfolioViewer(term.profile.Portfolios[0])
	term.portfolioViewer = portfolioViewer

	marketViewer := NewMarketViewer(term.market)
	term.marketViewer = marketViewer

	performanceViewer := NewPerformanceViewer(term.performance)
	term.performanceViewer = performanceViewer

	returnViewer := NewReturnViewer(term.performance)
	term.returnViewer = returnViewer

	term.setLayout()

	err := term.draw()
	if err != nil {
		return err
	}

	// This will lazily compute the performance and update the viewer
	go term.refreshPerformance()

	// Periodically refresh the market and portfolio data
	go term.doPeriodicRefresh()

	err = term.application.Run()
	if err != nil {
		return err
	}

	return nil
}

// Stop stops the terminal application
func (term *Terminal) Stop() {
	term.application.Stop()
}

func (term *Terminal) draw() error {
	err := term.market.Refresh()
	if err != nil {
		return err
	}

	err = term.profile.Portfolios[0].Refresh()
	if err != nil {
		return err
	}

	term.portfolioViewer.Draw()
	term.marketViewer.Draw()

	// The performance and return data have not been computed yet. However, that's handled by the viewer
	term.drawPerformance()

	return nil
}

func (term *Terminal) drawMarket() {
	term.marketViewer.Draw()
}

func (term *Terminal) drawPortfolio() {
	term.portfolioViewer.Draw()
}

func (term *Terminal) drawPerformance() {
	term.performanceViewer.Draw()
	term.returnViewer.Draw()
}

func (term *Terminal) refreshMarket() error {
	err := term.market.Refresh()
	if err != nil {
		return err
	}

	term.application.QueueUpdateDraw(func() {
		term.drawMarket()
	})

	return nil
}

func (term *Terminal) refreshPortfolio() error {
	err := term.profile.Portfolios[0].Refresh()
	if err != nil {
		return err
	}

	term.application.QueueUpdateDraw(func() {
		term.drawPortfolio()
	})

	return nil
}

func (term *Terminal) refreshPerformance() error {
	err := term.performance.Compute()
	if err != nil {
		return err
	}

	term.signalPerformance <- 0

	return nil
}

func (term *Terminal) doPeriodicRefresh() {
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ticker.C:
			term.refreshMarket()
			term.refreshPortfolio()

		case <-term.signalPerformance:
			term.application.QueueUpdateDraw(func() {
				term.drawPerformance()
			})
		}
	}
}

func (term *Terminal) setLayout() {
	grid := tview.NewGrid().SetRows(4, 0, 8, 7).SetColumns(0).SetBorders(false).
		AddItem(term.marketViewer.table, 0, 0, 1, 1, 0, 0, false).
		AddItem(term.portfolioViewer.table, 1, 0, 1, 1, 0, 0, false).
		AddItem(term.performanceViewer.table, 2, 0, 1, 1, 0, 0, false).
		AddItem(term.returnViewer.table, 3, 0, 1, 1, 0, 0, false)
	term.application.SetRoot(grid, true).SetInputCapture(term.keyCapture)
}

func (term *Terminal) keyCapture(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyRune && event.Rune() == 'q' {
		term.Stop()
		return nil
	}

	return event
}
