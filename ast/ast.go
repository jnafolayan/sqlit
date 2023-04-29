package ast

import (
	"fmt"
	"jnafolayan/sql-db/token"
	"strings"
)

type NodeType string

const (
	SELECT       NodeType = "SELECT"
	CREATE_TABLE NodeType = "CREATE_TABLE"
	INSERT       NodeType = "INSERT"

	INTEGER NodeType = "INTEGER"
	FLOAT   NodeType = "FLOAT"
	STRING  NodeType = "STRING"

	INFIX_EXPRESSION NodeType = "INFIX_EXPRESSION"
)

type Program struct {
	Statements []Statement
}

type Statement interface {
	Type() NodeType
	String() string
	statementNode()
}

type Expression interface {
	Type() NodeType
	String() string
	expressionNode()
}

type SelectStatement struct {
	Table     *token.Token
	Columns   []*token.Token
	Predicate Expression
}

func (ss *SelectStatement) statementNode() {}
func (ss *SelectStatement) Type() NodeType { return SELECT }
func (ss *SelectStatement) String() string {
	columns := []string{}
	for _, col := range ss.Columns {
		columns = append(columns, col.Literal)
	}

	cols := strings.Join(columns, ", ")

	predicate := ""
	if ss.Predicate != nil {
		predicate = fmt.Sprintf(" WHERE %s", ss.Predicate.String())
	}

	return fmt.Sprintf("SELECT %s FROM %s%s", cols, ss.Table.Literal, predicate)
}

type CreateTableStatement struct {
	Table   *token.Token
	Columns []*ColumnDefinition
}

func (cs *CreateTableStatement) statementNode() {}
func (cs *CreateTableStatement) Type() NodeType {
	return CREATE_TABLE
}
func (cs *CreateTableStatement) String() string {
	columns := []string{}
	for _, colDef := range cs.Columns {
		columns = append(columns, fmt.Sprintf("%s %s", colDef.Name.Literal, colDef.DataType.Literal))
	}

	cols := strings.Join(columns, ", ")
	return fmt.Sprintf("CREATE TABLE %s (%s)", cs.Table.Literal, cols)
}

type ColumnDefinition struct {
	Name     *token.Token
	DataType *token.Token
}

type InsertStatement struct {
	Table   *token.Token
	Columns []*token.Token
	Values  []Expression
}

func (is *InsertStatement) statementNode() {}
func (is *InsertStatement) Type() NodeType { return INSERT }
func (is *InsertStatement) String() string {
	columns := []string{}
	for _, col := range is.Columns {
		columns = append(columns, fmt.Sprintf("%s %s", col.Literal, col.Literal))
	}

	values := []string{}
	for _, val := range is.Values {
		values = append(values, val.String())
	}

	cols := strings.Join(columns, ", ")
	vals := strings.Join(values, ", ")
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", is.Table.Literal, cols, vals)
}

type IntegerLiteral struct {
	Token *token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) Type() NodeType  { return INTEGER }
func (il *IntegerLiteral) String() string {
	return fmt.Sprintf("%d", il.Value)
}

type FloatLiteral struct {
	Token *token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode() {}
func (fl *FloatLiteral) Type() NodeType  { return FLOAT }
func (fl *FloatLiteral) String() string {
	return fmt.Sprintf("%f", fl.Value)
}

type StringLiteral struct {
	Token *token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) Type() NodeType  { return STRING }
func (sl *StringLiteral) String() string {
	return sl.Value
}

type InfixExpression struct {
	Token    *token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) Type() NodeType  { return INFIX_EXPRESSION }
func (ie *InfixExpression) String() string {
	return fmt.Sprintf("%s %s %s", ie.Left.String(), ie.Operator, ie.Right.String())
}
