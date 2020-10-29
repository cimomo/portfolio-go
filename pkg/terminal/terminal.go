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
	application              *tview.Application
	root                     *tview.Pages
	profileFile              string
	profile                  *portfolio.Profile
	marketViewer             *MarketViewer
	profileViewer            *ProfileViewer
	profilePerformanceViewer *PerformanceViewer
	profileReturnViewer      *ReturnViewer
	portfolioViewers         []*PortfolioViewer
	performanceViewers       []*PerformanceViewer
	returnViewers            []*ReturnViewer
	currentViewer            int
	signalRedrawMarket       chan int
	signalRedrawProfile      chan int
	signalRedrawPortfolio    chan int
	signalRedrawPerformance  chan int
	signalSwitchViewer       chan int
}

// NewTerminal returns a new terminal window
func NewTerminal(profileFile string) *Terminal {
	return &Terminal{
		application:             tview.NewApplication(),
		profileFile:             profileFile,
		portfolioViewers:        make([]*PortfolioViewer, 0),
		performanceViewers:      make([]*PerformanceViewer, 0),
		returnViewers:           make([]*ReturnViewer, 0),
		currentViewer:           -1,
		signalRedrawMarket:      make(chan int),
		signalRedrawPortfolio:   make(chan int),
		signalRedrawProfile:     make(chan int),
		signalRedrawPerformance: make(chan int),
		signalSwitchViewer:      make(chan int),
	}
}

// Start starts the terminal application
func (term *Terminal) Start() error {
	profile, err := term.loadProfile("Main")
	if err != nil {
		return err
	}

	term.profile = profile

	term.setupViewers()

	term.initializeViewer()

	// Periodically refresh the market and portfolio data
	go term.doRefresh()

	err = term.application.Run()
	if err != nil {
		return err
	}

	return nil
}

func (term *Terminal) reload() error {
	profile, err := term.loadProfile("Main")
	if err != nil {
		return err
	}

	term.profile = profile

	term.marketViewer.Reload(profile.Market)
	term.profileViewer.Reload(profile)
	term.profilePerformanceViewer.Reload(profile.MergedPortfolio.Performance)
	term.profileReturnViewer.Reload(profile.MergedPortfolio.Performance)

	for i, portfolio := range profile.Portfolios {
		term.portfolioViewers[i].Reload(portfolio)
		term.performanceViewers[i].Reload(portfolio.Performance)
		term.returnViewers[i].Reload(portfolio.Performance)
	}

	err = term.switchViewer(term.currentViewer)
	if err != nil {
		return err
	}

	go term.computeAllPerformance()

	return nil
}

func (term *Terminal) loadProfile(name string) (*portfolio.Profile, error) {
	p := portfolio.NewProfile(name)

	err := p.Load(term.profileFile)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (term *Terminal) setupViewers() {
	marketViewer := NewMarketViewer(term.profile.Market)
	term.marketViewer = marketViewer

	profileViewer := NewProfileViewer(term.profile)
	term.profileViewer = profileViewer

	term.profilePerformanceViewer = NewPerformanceViewer(term.profile.MergedPortfolio.Performance)
	term.profileReturnViewer = NewReturnViewer(term.profile.MergedPortfolio.Performance)

	for _, portfolio := range term.profile.Portfolios {
		portfolioViewer := NewPortfolioViewer(portfolio)
		term.portfolioViewers = append(term.portfolioViewers, portfolioViewer)

		performanceViewer := NewPerformanceViewer(portfolio.Performance)
		term.performanceViewers = append(term.performanceViewers, performanceViewer)

		returnViewer := NewReturnViewer(portfolio.Performance)
		term.returnViewers = append(term.returnViewers, returnViewer)
	}
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

	err = term.profile.Refresh()
	if err != nil {
		return err
	}

	term.drawMarket()
	term.drawProfile()

	// The performance and return data may not have been computed yet. However, that's handled by the viewer
	term.drawPerformance(-1)

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
	term.initialize()

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

	term.signalSwitchViewer <- index

	return nil
}

func (term *Terminal) drawMarket() {
	term.marketViewer.Draw()
}

func (term *Terminal) drawProfile() {
	term.profileViewer.Draw()
}

func (term *Terminal) drawPortfolio(index int) error {
	if index < 0 || index >= len(term.portfolioViewers) {
		return errors.New("Viewer index out of range")
	}

	term.portfolioViewers[index].Draw()

	return nil
}

func (term *Terminal) drawPerformance(index int) error {
	if index >= len(term.performanceViewers) {
		return errors.New("Viewer index out of range")
	}

	if index < 0 {
		term.profilePerformanceViewer.Draw()
		term.profileReturnViewer.Draw()
	} else {
		term.performanceViewers[index].Draw()
		term.returnViewers[index].Draw()
	}

	return nil
}

func (term *Terminal) refreshMarket() error {
	err := term.profile.Market.Refresh()
	if err != nil {
		return err
	}

	term.signalRedrawMarket <- 0

	return nil
}

func (term *Terminal) refreshProfile() error {
	err := term.profile.Refresh()
	if err != nil {
		return err
	}

	term.signalRedrawProfile <- 0

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
	err := term.profile.MergedPortfolio.Performance.Compute()
	if err != nil {
		return err
	}

	for i := range term.profile.Portfolios {
		err = term.computePerformance(i)
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

	for {
		select {
		case <-ticker.C:
			go term.refreshMarket()

			if term.currentViewer < 0 {
				go term.refreshProfile()
			} else {
				go term.refreshPortfolio(term.currentViewer)
			}

		case <-term.signalRedrawMarket:
			term.application.QueueUpdateDraw(func() {
				term.drawMarket()
			})

		case <-term.signalRedrawProfile:
			term.application.QueueUpdateDraw(func() {
				term.drawProfile()
			})

		case <-term.signalRedrawPortfolio:
			term.application.QueueUpdateDraw(func() {
				term.drawPortfolio(term.currentViewer)
			})

		case <-term.signalRedrawPerformance:
			term.application.QueueUpdateDraw(func() {
				term.drawPerformance(term.currentViewer)
			})

		case index := <-term.signalSwitchViewer:
			term.currentViewer = index
			term.application.QueueUpdateDraw(func() {
				if index < 0 {
					term.root.SwitchToPage(term.profile.Name)
					term.drawHomepage()

				} else {
					term.root.SwitchToPage(term.portfolioViewers[index].portfolio.Name)
					term.drawPage(index)
				}
			})
		}
	}
}

func (term *Terminal) showHelp() {
	help := tview.NewModal().
		SetText("Help\n\n<h> This help menu\n<0>/<m> Home page\n<1>...<9> Switch to portfolio\n<r> Reload profile\n<q>/<Ctrl+c> Exit").
		AddButtons([]string{"Got it!"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			term.application.SetRoot(term.root, true)
		})
	help.SetTitle("Help")
	term.application.SetRoot(help, false)
}

func (term *Terminal) initialize() {
	pages := tview.NewPages()

	homepage := term.createHomepage()
	pages.AddPage(term.profile.Name, homepage, true, true)

	for i := range term.portfolioViewers {
		page := term.createPage(i)
		pages.AddPage(term.portfolioViewers[i].portfolio.Name, page, true, false)
	}

	term.application.SetRoot(pages, true).SetInputCapture(term.keyCapture)
	term.root = pages
}

func (term *Terminal) createHomepage() *tview.Grid {
	grid := tview.NewGrid().SetRows(4, 0, 8, 7).SetColumns(0).SetBorders(false)

	grid.AddItem(term.marketViewer.table, 0, 0, 1, 1, 0, 0, false).
		AddItem(term.profileViewer.table, 1, 0, 1, 1, 0, 0, false).
		AddItem(term.profilePerformanceViewer.table, 2, 0, 1, 1, 0, 0, false).
		AddItem(term.profileReturnViewer.table, 3, 0, 1, 1, 0, 0, false)

	return grid
}

func (term *Terminal) createPage(index int) *tview.Grid {
	grid := tview.NewGrid().SetRows(4, 0, 8, 7).SetColumns(0).SetBorders(false)

	grid.AddItem(term.marketViewer.table, 0, 0, 1, 1, 0, 0, false).
		AddItem(term.portfolioViewers[index].table, 1, 0, 1, 1, 0, 0, false).
		AddItem(term.performanceViewers[index].table, 2, 0, 1, 1, 0, 0, false).
		AddItem(term.returnViewers[index].table, 3, 0, 1, 1, 0, 0, false)

	return grid
}

func (term *Terminal) keyCapture(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyRune {
		rune := event.Rune()
		if rune == 'q' {
			term.Stop()
			return nil

		} else if rune >= '0' && rune <= '9' { // Note that '0' will open the main profile page
			index := int(rune - '1')
			term.switchViewer(index)
			return nil

		} else if rune == 'm' {
			term.switchViewer(-1)
			return nil

		} else if rune == 'h' {
			term.showHelp()
			return nil

		} else if rune == 'r' {
			term.reload()
			return nil
		}
	}

	return event
}
