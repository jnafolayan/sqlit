package engine

import (
	"bytes"
	"encoding/binary"
	"jnafolayan/sql-db/ast"
	"jnafolayan/sql-db/token"
	"strconv"
)

type memoryCell []byte

func (mc memoryCell) AsText() string {
	return string(mc)
}

func (mc memoryCell) AsInt() int64 {
	var i int64
	err := binary.Read(bytes.NewBuffer(mc), binary.BigEndian, &i)
	if err != nil {
		panic(err)
	}
	return i
}

type tableColumn struct {
	columnType ColumnType
	name       string
}

type table struct {
	columns []*tableColumn
	rows    [][]memoryCell
}

type MemoryBackend struct {
	tables map[string]*table
}

func NewMemoryBackend(existing map[string]*table) *MemoryBackend {
	tables := existing
	if existing == nil {
		tables = map[string]*table{}
	}

	return &MemoryBackend{
		tables: tables,
	}
}

func (mb *MemoryBackend) CreateTable(stmt *ast.CreateTableStatement) error {
	t := &table{}

	for _, col := range stmt.Columns {
		var colType ColumnType
		switch col.DataType.Type {
		case token.TEXT:
			colType = TEXT_COLUMN
		case token.INT:
			colType = INT_COLUMN
		default:
			return ErrInvalidDataType
		}

		t.columns = append(t.columns, &tableColumn{
			name:       col.Name.Literal,
			columnType: colType,
		})
	}

	mb.tables[stmt.Table.Literal] = t
	return nil
}

func (mb *MemoryBackend) Insert(stmt *ast.InsertStatement) error {
	t, ok := mb.tables[stmt.Table.Literal]
	if !ok {
		return ErrTableNotFound
	}

	// Generate colName -> colIndex map
	colNameToIdx := generateColNameToIndexMap(t.columns)

	// Allocate a row
	row := make([]memoryCell, len(t.columns))

	for i := range stmt.Columns {
		colName := stmt.Columns[i].Literal
		value := stmt.Values[i].String()

		colIdx, ok := colNameToIdx[colName]
		if !ok {
			return ErrColumnNotFound
		}

		var cellValue memoryCell

		switch t.columns[colIdx].columnType {
		case TEXT_COLUMN:
			cellValue = []byte(value)
		case INT_COLUMN:
			i, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return ErrInvalidDataType
			}

			buf := new(bytes.Buffer)
			err = binary.Write(buf, binary.BigEndian, i)
			if err != nil {
				panic(err)
			}

			cellValue = buf.Bytes()
		}

		row[colIdx] = cellValue
	}

	t.rows = append(t.rows, row)
	return nil
}

func (mb *MemoryBackend) Select(stmt *ast.SelectStatement) (*Result, error) {
	t, ok := mb.tables[stmt.Table.Literal]
	if !ok {
		return nil, ErrTableNotFound
	}

	resultRows := [][]Cell{}
	columns := []struct {
		Type ColumnType
		Name string
	}{}

	for _, col := range stmt.Columns {
		var colType ColumnType
		switch col.Type {
		case token.IDENTIFIER:
			colType = TEXT_COLUMN
		case token.INT:
			colType = INT_COLUMN
		default:
			return nil, ErrInvalidDataType
		}

		columns = append(columns, struct {
			Type ColumnType
			Name string
		}{
			Type: colType,
			Name: col.Literal,
		})
	}

	colNameToIdx := generateColNameToIndexMap(t.columns)
	for _, row := range t.rows {
		res := []Cell{}
		for _, col := range stmt.Columns {
			colIdx, ok := colNameToIdx[col.Literal]
			if !ok {
				return nil, ErrColumnNotFound
			}

			res = append(res, row[colIdx])
		}

		resultRows = append(resultRows, res)
	}

	return &Result{
		Rows:    resultRows,
		Columns: columns,
	}, nil
}

func generateColNameToIndexMap(columns []*tableColumn) map[string]int {
	colNameToIdx := map[string]int{}
	for i, col := range columns {
		colNameToIdx[col.name] = i
	}
	return colNameToIdx
}
