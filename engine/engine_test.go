package engine

import (
	"jnafolayan/sql-db/ast"
	"jnafolayan/sql-db/lexer"
	"jnafolayan/sql-db/parser"
	"testing"
)

func TestMemoryBackend(t *testing.T) {
	testStatement(t, "CREATE TABLE people (name TEXT, age INT, balance FLOAT)", func(tt *testing.T, i interface{}, err error) {
		if err != nil {
			t.Fatalf("error creating table: %s", err)
		}
	})
	testStatement(t, "INSERT INTO people (name, age, balance) VALUES ('John', 40, 0.5)", func(tt *testing.T, i interface{}, err error) {
		if err != nil {
			t.Fatalf("error inserting into table: %s", err)
		}
	})
	testStatement(t, "INSERT INTO people (name, age) VALUES ('Julia', 30)", func(tt *testing.T, i interface{}, err error) {
		if err != nil {
			t.Fatalf("error inserting into table: %s", err)
		}
	})

	testStatement(t, "SELECT name, age, balance FROM people", func(tt *testing.T, result interface{}, err error) {
		if err != nil {
			t.Fatalf("error selecting table: %s", err)
		}

		res, ok := result.(*FetchResult)
		if !ok {
			t.Fatalf("expected a table Result, got %T", result)
		}

		john := res.FetchAssoc()
		julia := res.FetchAssoc()

		if john["name"].AsText() != "John" {
			tt.Errorf("expected 'John', got %q", john["name"])
		}
		if john["age"].AsInt() != 40 {
			tt.Errorf("expected 40, got %d", john["age"])
		}
		if john["balance"].AsFloat() != 0.5 {
			tt.Errorf("expected 0.5, got %s", john["balance"].AsText())
		}

		if julia["name"].AsText() != "Julia" {
			tt.Errorf("expected 'Julia', got %q", julia["name"])
		}
		if julia["age"].AsInt() != 30 {
			tt.Errorf("expected 30, got %d", julia["age"])
		}
	})

	testStatement(t, "DELETE FROM people WHERE age=30", func(tt *testing.T, result interface{}, err error) {
		if err != nil {
			t.Fatalf("error deleting rows: %s", err)
		}

		res, ok := result.(*UpdateResult)
		if !ok {
			t.Fatalf("expected a UpdateResult, got %T", result)
		}

		if res.AffectedRows != 1 {
			tt.Errorf("expected 1 row to be deleted, got %d", res.AffectedRows)
		}
	})
	// Make sure it's deleted
	testStatement(t, "SELECT age FROM people WHERE age=30", func(tt *testing.T, result interface{}, err error) {
		if err != nil {
			t.Fatalf("error updating rows: %s", err)
		}

		res, ok := result.(*FetchResult)
		if !ok {
			tt.Fatalf("expected a FetchResult, got %T", result)
		}

		if len(res.Rows) != 0 {
			tt.Errorf("expected 0 rows, got %d", len(res.Rows))
		}
	})

	testStatement(t, "UPDATE people SET age=66 WHERE age=40", func(tt *testing.T, result interface{}, err error) {
		if err != nil {
			t.Fatalf("error updating rows: %s", err)
		}

		res, ok := result.(*UpdateResult)
		if !ok {
			t.Fatalf("expected a UpdateResult, got %T", result)
		}

		if res.AffectedRows != 1 {
			tt.Errorf("expected 1 row to be updated, got %d", res.AffectedRows)
		}
	})
	// Make sure it updated
	testStatement(t, "SELECT age FROM people WHERE age=66", func(tt *testing.T, result interface{}, err error) {
		if err != nil {
			t.Fatalf("error updating rows: %s", err)
		}

		res, ok := result.(*FetchResult)
		if !ok {
			tt.Fatalf("expected a FetchResult, got %T", result)
		}

		if len(res.Rows) != 1 {
			tt.Errorf("expected 1 row, got %d", len(res.Rows))
		}
	})

	testStatement(t, "DELETE FROM people", func(tt *testing.T, result interface{}, err error) {
		if err != nil {
			t.Fatalf("error deleting rows: %s", err)
		}

		res, ok := result.(*UpdateResult)
		if !ok {
			t.Fatalf("expected a UpdateResult, got %T", result)
		}

		if res.AffectedRows != 1 {
			tt.Errorf("expected 1 row to be deleted, got %d", res.AffectedRows)
		}
	})
}

var testTable = NewMemoryBackendTables()

func testStatement(t *testing.T, stmt string, callback func(*testing.T, interface{}, error)) {
	t.Run(stmt, func(tt *testing.T) {
		l := lexer.New(stmt)
		p := parser.New(l)
		program, err := p.Parse()
		if err != nil {
			t.Fatalf("error parsing statement: %s", err)
		}

		stmt := program.Statements[0]
		engine := NewMemoryBackend(testTable)

		switch st := stmt.(type) {
		case *ast.CreateTableStatement:
			callback(tt, nil, engine.CreateTable(st))
		case *ast.InsertStatement:
			callback(tt, nil, engine.Insert(st))
		case *ast.SelectStatement:
			res, err := engine.Select(st)
			callback(tt, res, err)
		case *ast.DeleteStatement:
			res, err := engine.Delete(st)
			callback(tt, res, err)
		case *ast.UpdateStatement:
			res, err := engine.Update(st)
			callback(tt, res, err)
		}
	})
}
