package terminal

import (
	"github.com/cimomo/portfolio-go/pkg/portfolio"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// ReturnViewer displays the trailing returns of a portfolio
type ReturnViewer struct {
	performance *portfolio.Performance
	table       *tview.Table
}

// NewReturnViewer returns a new viewer for the trailing returns of a portfolio
func NewReturnViewer(performance *portfolio.Performance) *ReturnViewer {
	return &ReturnViewer{
		performance: performance,
		table:       tview.NewTable().SetBorders(true),
	}
}

// Draw calculates the portfolio performance and refreshes the viewer
func (viewer *ReturnViewer) Draw() {
	viewer.drawHeader()
	viewer.drawPerformance()
}

func (viewer *ReturnViewer) drawHeader() {
	var cell *tview.TableCell
	header := []string{
		"Portfolio", "1-Month", "3-Month", "6-Month", "YTD",
		"1-Year", "3-Year", "5-Year", "10-Year", "Max",
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

func (viewer *ReturnViewer) drawPerformance() {
	if viewer.performance.StartDate.IsZero() {
		setString(viewer.table, "Computing ...", 1, 0, tcell.ColorWhite, tview.AlignLeft)
		setString(viewer.table, "Computing ...", 2, 0, tcell.ColorWhite, tview.AlignLeft)
		return
	}

	setString(viewer.table, viewer.performance.Portfolio.Name, 1, 0, tcell.ColorWhite, tview.AlignLeft)
	setPercentChange(viewer.table, viewer.performance.Result.Return.Max, 1, 10)

	setString(viewer.table, viewer.performance.Benchmark.Portfolio.Name, 2, 0, tcell.ColorWhite, tview.AlignLeft)
}
