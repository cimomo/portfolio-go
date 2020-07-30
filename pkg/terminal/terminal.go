package terminal

import (
	"errors"
	"time"

	"github.com/cimomo/portfolio-go/pkg/portfolio"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Terminal defines the main terminal window for portfolio visualization
type Terminal struct {
	application             *tview.Application
	root                    *tview.Grid
	profile                 *portfolio.Profile
	marketViewer            *MarketViewer
	profileViewer           *ProfileViewer
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

	profileViewer := NewProfileViewer(term.profile)
	term.profileViewer = profileViewer

	for _, portfolio := range term.profile.Portfolios {
		portfolioViewer := NewPortfolioViewer(portfolio)
		term.portfolioViewers = append(term.portfolioViewers, portfolioViewer)

		performanceViewer := NewPerformanceViewer(portfolio.Performance)
		term.performanceViewers = append(term.performanceViewers, performanceViewer)

		returnViewer := NewReturnViewer(portfolio.Performance)
		term.returnViewers = append(term.returnViewers, returnViewer)
	}

	term.initializeViewer()

	// Periodically refresh the market and portfolio data
	go term.doRefresh()

	err := term.application.Run()
	if err != nil {
		return err
	}

	return nil
}

// Stop stops the terminal application
func (term *Terminal) Stop() {
	term.application.Stop()
}

func (term *Terminal) drawHomepage() error {
	err := term.profile.Market.Refresh()
	if err != nil {
		return err
	}

	for _, portfolio := range term.profile.Portfolios {
		err = portfolio.Refresh()
		if err != nil {
			return err
		}
	}

	term.drawMarket()
	term.drawProfile()

	// The performance and return data may not have been computed yet. However, that's handled by the viewer
	term.drawPerformance(0)

	return nil
}

func (term *Terminal) drawPage(index int) error {
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

	// The performance and return data may not have been computed yet. However, that's handled by the viewer
	term.drawPerformance(index)

	return nil
}

func (term *Terminal) initializeViewer() error {
	term.initializeLayout()

	err := term.drawHomepage()
	if err != nil {
		return err
	}

	// This will lazily compute the performance and return data
	go term.computeAllPerformance()

	return nil
}

func (term *Terminal) switchViewer(index int) error {
	if index >= len(term.portfolioViewers) {
		return errors.New("Viewer index out of range")
	}

	if index < 0 {
		term.setLayoutForHomepage()

		err := term.drawHomepage()
		if err != nil {
			return err
		}

	} else {
		term.setLayoutForPage(index)

		err := term.drawPage(index)
		if err != nil {
			return err
		}
	}

	return nil
}

func (term *Terminal) drawMarket() {
	term.marketViewer.Draw()
}

func (term *Terminal) drawProfile() {
	term.profileViewer.Draw()
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

func (term *Terminal) computeAllPerformance() error {
	for i := range term.profile.Portfolios {
		err := term.computePerformance(i)
		if err != nil {
			return err
		}
	}

	// Sending a negative integer will cause the drawing thread to redraw the current viewer
	term.signalRedrawPerformance <- -1

	return nil
}

func (term *Terminal) computePerformance(index int) error {
	err := term.profile.Portfolios[index].Performance.Compute()
	if err != nil {
		return err
	}
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

		case i := <-term.signalRedrawPortfolio:
			if i >= 0 {
				index = i
			}

			term.application.QueueUpdateDraw(func() {
				term.drawPortfolio(index)
			})

		case i := <-term.signalRedrawPerformance:
			if i >= 0 {
				index = i
			}

			term.application.QueueUpdateDraw(func() {
				term.drawPerformance(index)
			})
		}
	}
}

func (term *Terminal) initializeLayout() {
	grid := tview.NewGrid().SetRows(4, 0, 8, 7).SetColumns(0).SetBorders(false)
	term.application.SetRoot(grid, true).SetInputCapture(term.keyCapture)
	term.root = grid
	term.setLayoutForHomepage()
}

func (term *Terminal) setLayoutForHomepage() {
	term.root.Clear()
	term.root.AddItem(term.marketViewer.table, 0, 0, 1, 1, 0, 0, false).
		AddItem(term.profileViewer.table, 1, 0, 1, 1, 0, 0, false).
		AddItem(term.performanceViewers[0].table, 2, 0, 1, 1, 0, 0, false).
		AddItem(term.returnViewers[0].table, 3, 0, 1, 1, 0, 0, false)
}

func (term *Terminal) setLayoutForPage(index int) {
	term.root.Clear()
	term.root.AddItem(term.marketViewer.table, 0, 0, 1, 1, 0, 0, false).
		AddItem(term.portfolioViewers[index].table, 1, 0, 1, 1, 0, 0, false).
		AddItem(term.performanceViewers[index].table, 2, 0, 1, 1, 0, 0, false).
		AddItem(term.returnViewers[index].table, 3, 0, 1, 1, 0, 0, false)
}

func (term *Terminal) keyCapture(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyRune {
		rune := event.Rune()
		if rune == 'q' {
			term.Stop()
			return nil

		} else if rune >= '1' && rune <= '9' {
			index := int(rune - '1')
			term.switchViewer(index)
			return nil

		} else if rune == 'h' {
			term.switchViewer(-1)
			return nil
		}
	}

	return event
}
