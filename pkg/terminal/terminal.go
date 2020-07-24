package terminal

import (
	"time"

	"github.com/cimomo/portfolio-go/pkg/portfolio"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Terminal defines the main terminal window for portfolio visualization
type Terminal struct {
	application             *tview.Application
	profile                 *portfolio.Profile
	marketViewer            *MarketViewer
	portfolioViewers        []*PortfolioViewer
	performanceViewers      []*PerformanceViewer
	returnViewers           []*ReturnViewer
	signalRedrawMarket      chan int
	signalRedrawPortfolio   chan int
	signalRedrawPerformance chan int
}

// NewTerminal returns a new terminal window
func NewTerminal(profile *portfolio.Profile) *Terminal {
	return &Terminal{
		application:             tview.NewApplication(),
		profile:                 profile,
		portfolioViewers:        make([]*PortfolioViewer, 0),
		performanceViewers:      make([]*PerformanceViewer, 0),
		returnViewers:           make([]*ReturnViewer, 0),
		signalRedrawMarket:      make(chan int),
		signalRedrawPortfolio:   make(chan int),
		signalRedrawPerformance: make(chan int),
	}
}

// Start starts the terminal application
func (term *Terminal) Start() error {
	marketViewer := NewMarketViewer(term.profile.Market)
	term.marketViewer = marketViewer

	for _, portfolio := range term.profile.Portfolios {
		portfolioViewer := NewPortfolioViewer(portfolio)
		term.portfolioViewers = append(term.portfolioViewers, portfolioViewer)

		performanceViewer := NewPerformanceViewer(portfolio.Performance)
		term.performanceViewers = append(term.performanceViewers, performanceViewer)

		returnViewer := NewReturnViewer(portfolio.Performance)
		term.returnViewers = append(term.returnViewers, returnViewer)
	}

	term.setLayout()

	err := term.draw(0)
	if err != nil {
		return err
	}

	// This will lazily compute the performance and update the viewer
	go term.refreshPerformance(0)

	// Periodically refresh the market and portfolio data
	go term.doRefresh()

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

func (term *Terminal) draw(index int) error {
	err := term.profile.Market.Refresh()
	if err != nil {
		return err
	}

	err = term.profile.Portfolios[index].Refresh()
	if err != nil {
		return err
	}

	term.drawMarket()
	term.drawPortfolio(index)

	// The performance and return data have not been computed yet. However, that's handled by the viewer
	term.drawPerformance(index)

	return nil
}

func (term *Terminal) drawMarket() {
	term.marketViewer.Draw()
}

func (term *Terminal) drawPortfolio(index int) {
	term.portfolioViewers[index].Draw()
}

func (term *Terminal) drawPerformance(index int) {
	term.performanceViewers[index].Draw()
	term.returnViewers[index].Draw()
}

func (term *Terminal) refreshMarket() error {
	err := term.profile.Market.Refresh()
	if err != nil {
		return err
	}

	term.signalRedrawMarket <- 0

	return nil
}

func (term *Terminal) refreshPortfolio(index int) error {
	err := term.profile.Portfolios[index].Refresh()
	if err != nil {
		return err
	}

	term.signalRedrawPortfolio <- index

	return nil
}

func (term *Terminal) refreshPerformance(index int) error {
	err := term.profile.Portfolios[index].Performance.Compute()
	if err != nil {
		return err
	}

	term.signalRedrawPerformance <- index

	return nil
}

func (term *Terminal) doRefresh() {
	ticker := time.NewTicker(time.Second * 10)
	index := 0

	for {
		select {
		case <-ticker.C:
			go term.refreshMarket()
			go term.refreshPortfolio(index)

		case <-term.signalRedrawMarket:
			term.application.QueueUpdateDraw(func() {
				term.drawMarket()
			})

		case index = <-term.signalRedrawPortfolio:
			term.application.QueueUpdateDraw(func() {
				term.drawPortfolio(index)
			})

		case index = <-term.signalRedrawPerformance:
			term.application.QueueUpdateDraw(func() {
				term.drawPerformance(index)
			})
		}
	}
}

func (term *Terminal) setLayout() {
	grid := tview.NewGrid().SetRows(4, 0, 8, 7).SetColumns(0).SetBorders(false).
		AddItem(term.marketViewer.table, 0, 0, 1, 1, 0, 0, false).
		AddItem(term.portfolioViewers[0].table, 1, 0, 1, 1, 0, 0, false).
		AddItem(term.performanceViewers[0].table, 2, 0, 1, 1, 0, 0, false).
		AddItem(term.returnViewers[0].table, 3, 0, 1, 1, 0, 0, false)
	term.application.SetRoot(grid, true).SetInputCapture(term.keyCapture)
}

func (term *Terminal) keyCapture(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyRune && event.Rune() == 'q' {
		term.Stop()
		return nil
	}

	return event
}
