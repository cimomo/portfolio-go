package terminal

import (
	"github.com/cimomo/portfolio-go/pkg/portfolio"
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
	viewer := NewPortfolioViewer(term)
	viewer.Connect()
}
