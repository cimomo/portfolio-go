package terminal

import (
	"github.com/cimomo/portfolio-go/pkg/portfolio"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// ProfileViewer displays real-time data for all portfolios in a profile
type ProfileViewer struct {
	profile *portfolio.Profile
	table   *tview.Table
}

// NewProfileViewer returns a new viewer for the real-time profile data
func NewProfileViewer(profile *portfolio.Profile) *ProfileViewer {
	return &ProfileViewer{
		profile: profile,
		table:   tview.NewTable().SetBorders(false),
	}
}

// Draw fetches the latest portfolio data and refreshes the viewer
func (viewer *ProfileViewer) Draw() {
	viewer.drawHeader()
	viewer.drawProfile()
}

func (viewer *ProfileViewer) drawHeader() {
	var cell *tview.TableCell
	header := []string{
		"NAME", "VALUE", "1-Day CHANGE%", "1-Day VALUE CHANGE$",
		"UNREALIZED GAIN/LOSS$", "UNREALIZED GAIN/LOSS%",
		"Allocation", "Target",
	}

	for c := 0; c < len(header); c++ {
		cell = tview.NewTableCell(header[c]).SetTextColor(tcell.ColorYellow).SetBackgroundColor(tcell.ColorDarkSlateGray).SetAttributes(tcell.AttrBold)
		if c < 2 {
			cell.SetAlign(tview.AlignLeft)
		} else {
			cell.SetAlign((tview.AlignRight))
		}
		viewer.table.SetCell(0, c, cell)
	}
}

func (viewer *ProfileViewer) drawProfile() {
	profile := viewer.profile
	portfolios := profile.Portfolios

	r := 1
	for _, portfolio := range portfolios {
		setString(viewer.table, portfolio.Name, r, 0, tcell.ColorWhite, tview.AlignLeft)
		setDollarAmount(viewer.table, portfolio.Status.Value, r, 1, tcell.ColorWhite)
		setPercentChange(viewer.table, portfolio.Status.RegularMarketChangePercent, r, 2)
		setDollarChange(viewer.table, portfolio.Status.RegularMarketChange, r, 3)
		setDollarChange(viewer.table, portfolio.Status.Unrealized, r, 4)
		setPercentChange(viewer.table, portfolio.Status.UnrealizedPercent, r, 5)
		setPercent(viewer.table, 0.0, r, 6, tcell.ColorWhite)
		setPercent(viewer.table, 0.0, r, 7, tcell.ColorWhite)

		r++
	}

	setString(viewer.table, "TOTAL", r, 0, tcell.ColorYellow, tview.AlignLeft)
	setPercent(viewer.table, 100.0, r, 6, tcell.ColorYellow)
	setPercent(viewer.table, 100.0, r, 7, tcell.ColorYellow)
}
