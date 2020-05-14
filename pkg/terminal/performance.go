package terminal

import (
	"github.com/cimomo/portfolio-go/pkg/portfolio"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// PerformanceViewer displays the historic performance of a portfolio
type PerformanceViewer struct {
	performance *portfolio.Performance
	table       *tview.Table
}

// NewPerformanceViewer returns a new viewer for the historic performance of a portfolio
func NewPerformanceViewer(performance *portfolio.Performance) *PerformanceViewer {
	return &PerformanceViewer{
		performance: performance,
		table:       tview.NewTable().SetBorders(true),
	}
}

// Draw calculates the portfolio performance and refreshes the viewer
func (viewer *PerformanceViewer) Draw() {
	viewer.drawHeader()
	viewer.drawPerformance()
}

func (viewer *PerformanceViewer) drawHeader() {
	var cell *tview.TableCell
	header := []string{
		"Portfolio", "Start Date", "Initial Balance", "Final Balance",
		"CAGR", "Stdev",
		"Best Year", "Worst Year", "Max Drawdown",
		"Sharpe Ratio", "US Market Correlation",
	}

	for c := 0; c < len(header); c++ {
		cell = tview.NewTableCell(header[c]).SetTextColor(tcell.ColorYellow).SetAttributes(tcell.AttrBold).SetExpansion(1)
		if c < 1 {
			cell.SetAlign(tview.AlignLeft)
		} else {
			cell.SetAlign((tview.AlignRight))
		}
		viewer.table.SetCell(0, c, cell)
	}
}

func (viewer *PerformanceViewer) drawPerformance() {
	if viewer.performance.StartDate.IsZero() {
		setString(viewer.table, "Computing ...", 1, 0, tcell.ColorWhite, tview.AlignLeft)
		return
	}

	setString(viewer.table, viewer.performance.Portfolio.Name, 1, 0, tcell.ColorWhite, tview.AlignLeft)
	setString(viewer.table, viewer.performance.StartDate.Format("2006-01-02"), 1, 1, tcell.ColorWhite, tview.AlignRight)
	setDollarAmount(viewer.table, viewer.performance.InitialBalance, 1, 2, tcell.ColorWhite)
	setDollarAmount(viewer.table, viewer.performance.FinalBalance, 1, 3, tcell.ColorWhite)
}
