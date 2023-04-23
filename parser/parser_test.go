package parser

import (
	"errors"
	"fmt"
	"jnafolayan/sql-db/ast"
	"jnafolayan/sql-db/lexer"
	"testing"
)

func TestParseSelectStatement(t *testing.T) {
	tests := []struct {
		input         string
		expectedError error
		expectedTable string
		expectedCols  []string
	}{
		{"SELECT * FROM people", nil, "people", []string{"*"}},
		{"SELECT name, age FROM people", nil, "people", []string{"name", "age"}},
		{
			"SELECT , FROM people",
			ErrEmptyColumnsList,
			"people",
			[]string{"name", "age"},
		},
	}

	for i, tt := range tests {
		testName := fmt.Sprintf("SELECT_%d", i)
		t.Run(testName, func(sub *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			p.OmitErrorLocation = true
			program, err := p.Parse()

			if err != nil {
				if tt.expectedError == nil {
					sub.Fatalf("expected no error, got %q", err)
				}
				if !errors.Is(err, tt.expectedError) {
					sub.Fatalf("expected %q error, got %q", tt.expectedError, err)
				}
				return
			} else if tt.expectedError != nil {
				sub.Fatalf("expected %v error, got no error", err)
			}

			if len(program.Statements) != 1 {
				sub.Fatalf("expected 1 statement, got %d", len(program.Statements))
			}

			stmt := program.Statements[0]
			if stmt.Type() != ast.SELECT {
				t.Fatalf("expected select statement, got %s", stmt.Type())
			}

			selectStmt := stmt.(*ast.SelectStatement)
			if selectStmt.Table.Literal != tt.expectedTable {
				t.Fatalf("expected table %q, got %q", tt.expectedTable, selectStmt.Table.Literal)
			}

			for i, col := range selectStmt.Columns {
				if col.Literal != tt.expectedCols[i] {
					t.Errorf("expected %q column, got %q", tt.expectedCols[i], col.Literal)
				}
			}
		})
	}
}

func TestParseCreateTableStatement(t *testing.T) {
	type colDef struct {
		name    string
		colType string
	}
	tests := []struct {
		input         string
		expectedError error
		expectedTable string
		expectedCols  []colDef
	}{
		{"CREATE TABLE people;", nil, "people", []colDef{}},
		{
			"CREATE TABLE people (name TEXT, age INT)",
			nil,
			"people",
			[]colDef{
				{"name", "TEXT"},
				{"age", "INT"},
			},
		},
	}

	for i, tt := range tests {
		testName := fmt.Sprintf("CREATE_%d", i)
		t.Run(testName, func(sub *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			p.OmitErrorLocation = true
			program, err := p.Parse()

			if err != nil {
				if tt.expectedError == nil {
					sub.Fatalf("expected no error, got %q", err)
				}
				if !errors.Is(err, tt.expectedError) {
					sub.Fatalf("expected %q error, got %q", tt.expectedError, err)
				}
				return
			} else if tt.expectedError != nil {
				sub.Fatalf("expected %v error, got no error", err)
			}

			if len(program.Statements) != 1 {
				sub.Fatalf("expected 1 statement, got %d", len(program.Statements))
			}

			stmt := program.Statements[0]
			if stmt.Type() != ast.CREATE_TABLE {
				t.Fatalf("expected create table statement, got %s", stmt.Type())
			}

			createStmt := stmt.(*ast.CreateTableStatement)
			if createStmt.Table.Literal != tt.expectedTable {
				t.Fatalf("expected table %q, got %q", tt.expectedTable, createStmt.Table.Literal)
			}

			for i, col := range createStmt.Columns {
				if col.Name.Literal != tt.expectedCols[i].name {
					t.Errorf("expected %q column, got %q", tt.expectedCols[i].name, col.Name.Literal)
					continue
				}
				if col.DataType.Literal != tt.expectedCols[i].colType {
					t.Errorf("expected %q column, got %q", tt.expectedCols[i].colType, col.DataType.Literal)
				}
			}
		})
	}
}

func TestParseInsertStatement(t *testing.T) {
	tests := []struct {
		input          string
		expectedError  error
		expectedTable  string
		expectedCols   []string
		expectedValues []string
	}{
		{"INSERT INTO people (name) VALUES ('jake')", nil, "people", []string{"name"}, []string{"jake"}},
	}

	for i, tt := range tests {
		testName := fmt.Sprintf("INSERT_%d", i)
		t.Run(testName, func(sub *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			p.OmitErrorLocation = true
			program, err := p.Parse()

			if err != nil {
				if tt.expectedError == nil {
					sub.Fatalf("expected no error, got %q", err)
				}
				if !errors.Is(err, tt.expectedError) {
					sub.Fatalf("expected %q error, got %q", tt.expectedError, err)
				}
				return
			} else if tt.expectedError != nil {
				sub.Fatalf("expected %v error, got no error", err)
			}

			if len(program.Statements) != 1 {
				sub.Fatalf("expected 1 statement, got %d", len(program.Statements))
			}

			stmt := program.Statements[0]
			if stmt.Type() != ast.INSERT {
				t.Fatalf("expected insert statement, got %s", stmt.Type())
			}

			insertStmt := stmt.(*ast.InsertStatement)
			if insertStmt.Table.Literal != tt.expectedTable {
				t.Fatalf("expected table %q, got %q", tt.expectedTable, insertStmt.Table.Literal)
			}

			for i, col := range insertStmt.Columns {
				if col.Literal != tt.expectedCols[i] {
					t.Errorf("expected %q column, got %q", tt.expectedCols[i], col.Literal)
				}
			}

			for i, val := range insertStmt.Values {
				if val.String() != tt.expectedValues[i] {
					t.Errorf("expected %q column, got %q", tt.expectedCols[i], val.String())
				}
			}
		})
	}
}
