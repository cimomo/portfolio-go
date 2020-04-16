package terminal

import (
	"time"

	"github.com/cimomo/portfolio-go/pkg/portfolio"
	"github.com/rivo/tview"
)

// Terminal defines the main terminal window for portfolio visualization
type Terminal struct {
	application     *tview.Application
	portfolio       *portfolio.Portfolio
	portfolioViewer *PortfolioViewer
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
	viewer := NewPortfolioViewer(term)
	term.portfolioViewer = viewer
	viewer.Connect()
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
