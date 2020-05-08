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
	viewer := PerformanceViewer{
		performance: performance,
		table:       tview.NewTable().SetBorders(true),
	}

	viewer.drawHeader()

	return &viewer
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
