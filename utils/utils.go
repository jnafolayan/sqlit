package utils

import (
	"fmt"
	"jnafolayan/sql-db/engine"
	"math"
	"strings"
)

func FormatSelectResult(result *engine.FetchResult) string {
	cellSizes := map[int]int{}
	for i := range result.Columns {
		cellSizes[i] = getLargestCellSize(i, result) + 2
	}

	// print header
	var header strings.Builder
	for i, col := range result.Columns {
		if i == 0 {
			header.WriteString("|")
		}
		header.WriteString(alignText(col.Name, cellSizes[i], " "))
		header.WriteString("|")
	}

	underline := strings.Repeat("=", header.Len()+5)
	header.WriteString("\n" + underline)

	var rowsBuilder strings.Builder

	for _, row := range result.Rows {
		var rowBuilder strings.Builder
		for i, cell := range row {
			resCol := result.Columns[i]
			content := ""
			if resCol.Type == engine.INT_COLUMN {
				content = fmt.Sprintf("%d", cell.AsInt())
			} else if resCol.Type == engine.FLOAT_COLUMN {
				content = fmt.Sprintf("%f", cell.AsFloat())
			} else {
				content = cell.AsText()
			}

			if i == 0 {
				rowBuilder.WriteString("|")
			}
			rowBuilder.WriteString(alignText(content, cellSizes[i], " "))
			rowBuilder.WriteString("|")
		}
		rowsBuilder.WriteString(rowBuilder.String())
		rowsBuilder.WriteString("\n")
	}

	return fmt.Sprintf("%s\n%s", header.String(), rowsBuilder.String())
}

func alignText(str string, length int, prefix string) string {
	res := str
	if len(res) < length {
		res = fmt.Sprintf(" %s%s", res, strings.Repeat(prefix, length-len(res)))
	}
	return res
}

func getLargestCellSize(column int, result *engine.FetchResult) int {
	largest := 0.
	for _, row := range result.Rows {
		content := ""
		resCol := result.Columns[column]
		if resCol.Type == engine.INT_COLUMN {
			content = fmt.Sprintf("%d", row[column].AsInt())
		} else if resCol.Type == engine.FLOAT_COLUMN {
			content = fmt.Sprintf("%f", row[column].AsFloat())
		} else {
			content = row[column].AsText()
		}
		largest = math.Max(largest, float64(len(content)))
	}
	return int(largest)
}
