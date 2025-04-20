package pkg

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"strings"
)

func RenderDiffs(oldColoured string, newColoured string) {

	oldLines := strings.Split(oldColoured, "\n")
	newLines := strings.Split(newColoured, "\n")
	maxLines := max(len(oldLines), len(newLines))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Old text:", "New text:"})
	table.SetAutoWrapText(false)
	table.SetReflowDuringAutoWrap(false)
	table.SetBorder(false)
	table.SetRowLine(true)
	table.SetColumnSeparator(" | ")
	table.SetColWidth(50)

	for i := range maxLines {
		var o, n string
		if i < len(oldLines) {
			o = oldLines[i]
		}
		if i < len(newLines) {
			n = newLines[i]
		}
		row := []string{o, n}
		fmt.Printf("row: %v\n", row)
		table.Append(row)
	}

	table.Render()
}
