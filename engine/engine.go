package engine

import (
	"errors"
	"jnafolayan/sql-db/ast"
	"jnafolayan/sql-db/lib"
)

type ColumnType string

const (
	INT_COLUMN   ColumnType = "INT"
	FLOAT_COLUMN ColumnType = "FLOAT"
	TEXT_COLUMN  ColumnType = "TEXT"
)

type Cell interface {
	AsText() string
	AsInt() int64
	AsFloat() float64
}

type RowAssoc map[string]Cell

type ResultColumn struct {
	Type ColumnType
	Name string
}

type FetchResult struct {
	Rows    [][]Cell
	Columns []*ResultColumn

	it *lib.VirtualIterator[RowAssoc]
}

func (r *FetchResult) FetchAssoc() RowAssoc {
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

type DeleteResult struct {
	affectedRows int
}

var (
	ErrInvalidDataType = errors.New("Invalid datatype")
	ErrTableNotFound   = errors.New("Table not found")
	ErrTableExists     = errors.New("Table already exists")
	ErrColumnNotFound  = errors.New("Column not found")
)

type Engine interface {
	Select(*ast.SelectStatement) (*FetchResult, error)
	CreateTable(*ast.CreateTableStatement) error
	Insert(*ast.InsertStatement) error
	Delete(*ast.DeleteStatement) error
}
