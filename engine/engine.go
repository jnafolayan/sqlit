package engine

import (
	"errors"
	"jnafolayan/sql-db/ast"
	"jnafolayan/sql-db/lib"
)

type ColumnType string

const (
	INT_COLUMN  ColumnType = "INT"
	TEXT_COLUMN ColumnType = "TEXT"
)

type Cell interface {
	AsText() string
	AsInt() int64
}

type RowAssoc map[string]Cell

type Result struct {
	Rows    [][]Cell
	Columns []struct {
		Type ColumnType
		Name string
	}

	it *lib.VirtualIterator[RowAssoc]
}

func (r *Result) FetchAssoc() RowAssoc {
	if r.it == nil {
		// TODO: optimize 'res'
		r.it = lib.NewVirtualIterator(len(r.Rows), func(cursor int) RowAssoc {
			res := RowAssoc{}
			for i, col := range r.Columns {
				res[col.Name] = r.Rows[cursor][i]
			}
			return res
		})
	}

	return r.it.Next()
}

var (
	ErrInvalidDataType = errors.New("Invalid datatype")
	ErrTableNotFound   = errors.New("Table not found")
	ErrColumnNotFound  = errors.New("Column not found")
)

type Engine interface {
	Select(*ast.SelectStatement) (*Result, error)
	CreateTable(*ast.CreateTableStatement) error
	Insert(*ast.InsertStatement) error
}
