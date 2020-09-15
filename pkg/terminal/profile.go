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

// Reload updates the profile data object
func (viewer *ProfileViewer) Reload(profile *portfolio.Profile) {
	viewer.profile = profile
}

// Draw fetches the latest portfolio data and refreshes the viewer
func (viewer *ProfileViewer) Draw() {
	viewer.table.Clear()
	viewer.drawHeader()
	viewer.drawProfile()
}

func (viewer *ProfileViewer) drawHeader() {
	var cell *tview.TableCell
	header := []string{
		"NAME", "COST BASIS", "VALUE", "1-DAY CHANGE%", "1-DAY VALUE CHANGE$",
		"UNREALIZED GAIN/LOSS$", "UNREALIZED GAIN/LOSS%",
		"ALLOCATION", "TARGET",
	}

	for c := 0; c < len(header); c++ {
		cell = tview.NewTableCell(header[c]).SetTextColor(tcell.ColorYellow).SetBackgroundColor(tcell.ColorDarkSlateGray).SetAttributes(tcell.AttrBold)
		if c < 1 {
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
		setDollarAmount(viewer.table, portfolio.CostBasis, r, 1, tcell.ColorWhite)
		setDollarAmount(viewer.table, portfolio.Status.Value, r, 2, tcell.ColorWhite)
		setPercentChange(viewer.table, portfolio.Status.RegularMarketChangePercent, r, 3)
		setDollarChange(viewer.table, portfolio.Status.RegularMarketChange, r, 4)
		setDollarChange(viewer.table, portfolio.Status.Unrealized, r, 5)
		setPercentChange(viewer.table, portfolio.Status.UnrealizedPercent, r, 6)
		setPercent(viewer.table, profile.Status.Allocation[portfolio.Name], r, 7, tcell.ColorWhite)
		setPercent(viewer.table, profile.TargetAllocation[portfolio.Name], r, 8, tcell.ColorWhite)

		r++
	}

	setString(viewer.table, "Cash", r, 0, tcell.ColorWhite, tview.AlignLeft)
	setDollarAmount(viewer.table, profile.Cash, r, 1, tcell.ColorWhite)
	setDollarAmount(viewer.table, profile.Cash, r, 2, tcell.ColorWhite)
	setPercentChange(viewer.table, 0, r, 3)
	setDollarChange(viewer.table, 0, r, 4)
	setDollarChange(viewer.table, 0, r, 5)
	setPercentChange(viewer.table, 0, r, 6)
	setPercent(viewer.table, profile.Status.Allocation["cash"], r, 7, tcell.ColorWhite)
	setPercent(viewer.table, profile.TargetAllocation["cash"], r, 8, tcell.ColorWhite)

	r++

	setString(viewer.table, "TOTAL", r, 0, tcell.ColorYellow, tview.AlignLeft)
	setDollarAmount(viewer.table, profile.CostBasis, r, 1, tcell.ColorYellow)
	setDollarAmount(viewer.table, profile.Status.Value, r, 2, tcell.ColorYellow)
	setPercentChange(viewer.table, profile.Status.RegularMarketChangePercent, r, 3)
	setDollarChange(viewer.table, profile.Status.RegularMarketChange, r, 4)
	setDollarChange(viewer.table, profile.Status.Unrealized, r, 5)
	setPercentChange(viewer.table, profile.Status.UnrealizedPercent, r, 6)
	setPercent(viewer.table, 100.0, r, 7, tcell.ColorYellow)
	setPercent(viewer.table, 100.0, r, 8, tcell.ColorYellow)
}
