package repl

import (
	"bufio"
	"fmt"
	"io"
	"jnafolayan/sql-db/ast"
	"jnafolayan/sql-db/engine"
	"jnafolayan/sql-db/lexer"
	"jnafolayan/sql-db/parser"
	"math"
	"os"
	"strings"
	"time"
)

const PROMPT = "# "

func Start(input io.Reader, output io.Writer) {
	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanLines)

	backend := engine.NewMemoryBackend(nil)
	fmt.Println("SQLit version 1.0")

	for {
		fmt.Fprintf(output, PROMPT)
		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program, err := p.Parse()
		if err != nil {
			fmt.Fprintf(os.Stderr, "program error: %s\n", err)
			continue
		}

		end := len(program.Statements) - 1

		startTime := time.Now()
		for i, stmt := range program.Statements {
			switch st := stmt.(type) {
			case *ast.CreateTableStatement:
				err := backend.CreateTable(st)
				if err != nil {
					fmt.Fprintf(os.Stderr, "program error: %s\n", err)
				} else if i == end {
					fmt.Fprintln(output, "ok")
				}
			case *ast.InsertStatement:
				err := backend.Insert(st)
				if err != nil {
					fmt.Fprintf(os.Stderr, "program error: %s\n", err)
				} else if i == end {
					fmt.Fprintln(output, "ok")
				}
			case *ast.SelectStatement:
				res, err := backend.Select(st)
				if err != nil {
					fmt.Fprintf(os.Stderr, "program error: %s\n", err)
				} else if i == end {
					printSelectResult(output, res)
					fmt.Fprintln(output, "ok")
				}
			}
			if i == end {
				duration := time.Now().Sub(startTime).Seconds()
				fmt.Fprintf(output, "ok (took %.2fs)\n", duration)
			}
		}
	}
}

func printSelectResult(output io.Writer, result *engine.Result) {
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

	fmt.Fprintln(output, header.String())
	fmt.Fprintln(output, underline)

	for _, row := range result.Rows {
		var rowBuilder strings.Builder
		for i, cell := range row {
			resCol := result.Columns[i]
			content := ""
			if resCol.Type == engine.INT_COLUMN {
				content = fmt.Sprintf("%d", cell.AsInt())
			} else {
				content = cell.AsText()
			}

			if i == 0 {
				rowBuilder.WriteString("|")
			}
			rowBuilder.WriteString(alignText(content, cellSizes[i], " "))
			rowBuilder.WriteString("|")
		}
		fmt.Fprintln(output, rowBuilder.String())
	}
}

func alignText(str string, length int, prefix string) string {
	res := str
	if len(res) < length {
		res = fmt.Sprintf(" %s%s", res, strings.Repeat(prefix, length-len(res)))
	}
	return res
}

func getLargestCellSize(column int, result *engine.Result) int {
	largest := 0.
	for _, row := range result.Rows {
		content := ""
		resCol := result.Columns[column]
		if resCol.Type == engine.INT_COLUMN {
			content = fmt.Sprintf("%d", row[column].AsInt())
		} else {
			content = row[column].AsText()
		}
		largest = math.Max(largest, float64(len(content)))
	}
	return int(largest)
}
