package terminal

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// HelpViewer displays the help message
type HelpViewer struct {
	help  map[string]string
	table *tview.Table
}

// NewHelpViewer returns a new viewer for help message
func NewHelpViewer() *HelpViewer {
	return &HelpViewer{
		help: map[string]string{
			"<h>":            "This help message",
			"<0>/<m>":        "Switch to home page",
			"<1>...<9>":      "Switch to portfolio",
			"<r>":            "Reload profile",
			"<q>/<Ctrl>+<c>": "Exit",
		},
		table: tview.NewTable().SetBorders(false),
	}
}

// Draw populates the help message
func (viewer *HelpViewer) Draw() {
	viewer.table.Clear()

	viewer.table.SetTitle("Help").SetBorder(true).SetBorderPadding(1, 1, 2, 2).SetBackgroundColor(tcell.ColorBlue)

	r := 0
	for k, v := range viewer.help {
		setString(viewer.table, k, r, 0, tcell.ColorWhite, tview.AlignLeft)
		setString(viewer.table, v, r, 1, tcell.ColorWhite, tview.AlignLeft)
		r++
	}
}
