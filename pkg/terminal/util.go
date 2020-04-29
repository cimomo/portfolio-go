package terminal

import (
	"math"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func setPercentChange(table *tview.Table, value float64, r int, c int) {
	color := tcell.ColorGreen
	if value < 0 {
		color = tcell.ColorRed
	}
	setPercent(table, value, r, c, color)
}

func setDollarChange(table *tview.Table, value float64, r int, c int) {
	color := tcell.ColorGreen
	if value < 0 {
		color = tcell.ColorRed
	}
	setDollarAmount(table, value, r, c, color)
}

func setPercent(table *tview.Table, value float64, r int, c int, color tcell.Color) {
	setFloat64(table, value, "%.2f%%", r, c, color, tview.AlignRight)
}

func setDollarAmount(table *tview.Table, value float64, r int, c int, color tcell.Color) {
	formatter := "$%.2f"
	if value < 0 {
		formatter = "-$%.2f"
	}

	v := math.Abs(value)
	setFloat64(table, v, formatter, r, c, color, tview.AlignRight)
}

func setQuantity(table *tview.Table, value float64, r int, c int, align int) {
	setFloat64(table, value, "%.2f", r, c, tcell.ColorWhite, align)
}

func setFloat64(table *tview.Table, value float64, formatter string, r int, c int, color tcell.Color, align int) {
	printer := message.NewPrinter(language.English)
	fValue := printer.Sprintf(formatter, value)
	setString(table, fValue, r, c, color, align)
}

func setString(table *tview.Table, value string, r int, c int, color tcell.Color, align int) {
	cell := tview.NewTableCell(value).SetTextColor(color).SetAlign(align).SetExpansion(1)
	table.SetCell(r, c, cell)
}
