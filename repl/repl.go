package repl

import (
	"bufio"
	"fmt"
	"io"
	"jnafolayan/sql-db/ast"
	"jnafolayan/sql-db/engine"
	"jnafolayan/sql-db/lexer"
	"jnafolayan/sql-db/parser"
	"jnafolayan/sql-db/utils"
	"os"
	"time"
)

const PROMPT = "sqlit> "

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
	loop:
		for i, stmt := range program.Statements {
			switch st := stmt.(type) {
			case *ast.CreateTableStatement:
				err := backend.CreateTable(st)
				if err != nil {
					fmt.Fprintf(os.Stderr, "program error: %s\n", err)
					break loop
				}
			case *ast.InsertStatement:
				err := backend.Insert(st)
				if err != nil {
					fmt.Fprintf(os.Stderr, "program error: %s\n", err)
					break loop
				}
			case *ast.SelectStatement:
				res, err := backend.Select(st)
				if err != nil {
					fmt.Fprintf(os.Stderr, "program error: %s\n", err)
					break loop
				} else if i == end {
					// Print only if result is not empty
					if len(res.Rows) != 0 {
						fmt.Fprintln(output, utils.FormatSelectResult(res))
					}
				}
			case *ast.DeleteStatement:
				res, err := backend.Delete(st)
				if err != nil {
					fmt.Fprintf(os.Stderr, "program error: %s\n", err)
					break loop
				} else if i == end {
					// Print only if result is not empty
					fmt.Fprintf(output, "affected rows: %d\n", res.AffectedRows)
				}
			}
			if i == end {
				duration := time.Now().Sub(startTime).Seconds()
				fmt.Fprintf(output, "ok (took %.2fs)\n", duration)
			}
		}
	}
}
