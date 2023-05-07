package main

import (
	"fmt"
	"jnafolayan/sql-db/ast"
	"jnafolayan/sql-db/engine"
	"jnafolayan/sql-db/lexer"
	"jnafolayan/sql-db/parser"
	"jnafolayan/sql-db/utils"
	"strings"
	"time"
)

var tables engine.MemoryTables

func init() {
	tables = engine.NewMemoryBackendTables()
}

//export logText
func logText(message string)

func main() {
	// populate database
	query := `CREATE TABLE people (name TEXT, age INT, address TEXT);
INSERT INTO people (name, age, address) VALUES ('John Doe', 43, 'Earth, Solar');
INSERT INTO people (name, age, address) VALUES ('Tracy White', 22, 'Venus, Solar');
INSERT INTO people (name, age, address) VALUES ('Tom Fischer', 26, 'Mercury, Solar');`
	program, err := parse(query)
	if err != nil {
		logText(err.Error())
	} else {
		execute(program)
	}
}

//export run
func run(input string) {
	program, err := parse(input)
	if err != nil {
		logText(err.Error())
	} else {
		logText(execute(program))
	}
}

//export parse
func parse(input string) (*ast.Program, error) {
	l := lexer.New(input)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		return nil, fmt.Errorf("program error: %s", err)
	}

	return program, nil
}

//export execute
func execute(program *ast.Program) string {
	var result strings.Builder

	backend := engine.NewMemoryBackend(tables)
	end := len(program.Statements) - 1

	startTime := time.Now()
loop:
	for i, stmt := range program.Statements {
		switch st := stmt.(type) {
		case *ast.CreateTableStatement:
			err := backend.CreateTable(st)
			if err != nil {
				result.WriteString(fmt.Errorf("program error: %s\n", err).Error())
				break loop
			}
		case *ast.InsertStatement:
			err := backend.Insert(st)
			if err != nil {
				result.WriteString(fmt.Errorf("program error: %s\n", err).Error())
				break loop
			}
		case *ast.SelectStatement:
			res, err := backend.Select(st)
			if err != nil {
				result.WriteString(fmt.Errorf("program error: %s\n", err).Error())
				break loop
			} else if i == end {
				result.WriteString(utils.FormatSelectResult(res))
			}
		case *ast.DeleteStatement:
			res, err := backend.Delete(st)
			if err != nil {
				result.WriteString(fmt.Errorf("program error: %s\n", err).Error())
				break loop
			} else if i == end {
				// Print only if result is not empty
				result.WriteString(fmt.Sprintf("affected rows: %d\n", res.AffectedRows))
			}
		case *ast.UpdateStatement:
			res, err := backend.Update(st)
			if err != nil {
				result.WriteString(fmt.Errorf("program error: %s\n", err).Error())
				break loop
			} else if i == end {
				// Print only if result is not empty
				result.WriteString(fmt.Sprintf("affected rows: %d\n", res.AffectedRows))
			}

		}

		if i == end {
			duration := time.Now().Sub(startTime).Seconds()
			result.WriteString(fmt.Sprintf("ok (took %.2fs)\n", duration))
		}
	}

	return result.String()
}
