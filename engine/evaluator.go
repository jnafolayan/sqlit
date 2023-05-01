package engine

import (
	"errors"
	"fmt"
	"jnafolayan/sql-db/ast"
	"strings"
)

type infixEvalFn func(ast.Expression, string, ast.Expression) (ast.Expression, error)

var infixEvalFns map[string]infixEvalFn

func init() {
	infixEvalFns = map[string]infixEvalFn{}

	// INTEGER + INTEGER
	registerInfixEvalFn(ast.INTEGER, "+", ast.INTEGER, func(e1 ast.Expression, s string, e2 ast.Expression) (ast.Expression, error) {
		a, _ := e1.(*ast.IntegerLiteral)
		b, _ := e2.(*ast.IntegerLiteral)
		return &ast.IntegerLiteral{
			Value: a.Value + b.Value,
		}, nil
	})

	// INTEGER - INTEGER
	registerInfixEvalFn(ast.INTEGER, "-", ast.INTEGER, func(e1 ast.Expression, s string, e2 ast.Expression) (ast.Expression, error) {
		a, _ := e1.(*ast.IntegerLiteral)
		b, _ := e2.(*ast.IntegerLiteral)
		return &ast.IntegerLiteral{
			Value: a.Value - b.Value,
		}, nil
	})

	// INTEGER = INTEGER
	registerInfixEvalFn(ast.INTEGER, "=", ast.INTEGER, func(e1 ast.Expression, s string, e2 ast.Expression) (ast.Expression, error) {
		a, _ := e1.(*ast.IntegerLiteral)
		b, _ := e2.(*ast.IntegerLiteral)
		return &ast.Boolean{
			Value: a.Value == b.Value,
		}, nil
	})

	// FLOAT + FLOAT
	registerInfixEvalFn(ast.FLOAT, "+", ast.FLOAT, func(e1 ast.Expression, s string, e2 ast.Expression) (ast.Expression, error) {
		a, _ := e1.(*ast.FloatLiteral)
		b, _ := e2.(*ast.FloatLiteral)
		return &ast.FloatLiteral{
			Value: a.Value + b.Value,
		}, nil
	})

	// FLOAT - FLOAT
	registerInfixEvalFn(ast.FLOAT, "-", ast.FLOAT, func(e1 ast.Expression, s string, e2 ast.Expression) (ast.Expression, error) {
		a, _ := e1.(*ast.FloatLiteral)
		b, _ := e2.(*ast.FloatLiteral)
		return &ast.FloatLiteral{
			Value: a.Value - b.Value,
		}, nil
	})

	// FLOAT = FLOAT
	registerInfixEvalFn(ast.FLOAT, "=", ast.FLOAT, func(e1 ast.Expression, s string, e2 ast.Expression) (ast.Expression, error) {
		a, _ := e1.(*ast.FloatLiteral)
		b, _ := e2.(*ast.FloatLiteral)
		return &ast.Boolean{
			Value: a.Value == b.Value,
		}, nil
	})

	// STRING = STRING
	registerInfixEvalFn(ast.STRING, "=", ast.STRING, func(e1 ast.Expression, s string, e2 ast.Expression) (ast.Expression, error) {
		a, _ := e1.(*ast.StringLiteral)
		b, _ := e2.(*ast.StringLiteral)
		return &ast.Boolean{
			Value: a.Value == b.Value,
		}, nil
	})

	// INTEGER != INTEGER
	registerInfixEvalFn(ast.INTEGER, "!=", ast.INTEGER, func(e1 ast.Expression, s string, e2 ast.Expression) (ast.Expression, error) {
		a, _ := e1.(*ast.IntegerLiteral)
		b, _ := e2.(*ast.IntegerLiteral)
		return &ast.Boolean{
			Value: a.Value != b.Value,
		}, nil
	})

	// FLOAT != FLOAT
	registerInfixEvalFn(ast.FLOAT, "!=", ast.FLOAT, func(e1 ast.Expression, s string, e2 ast.Expression) (ast.Expression, error) {
		a, _ := e1.(*ast.FloatLiteral)
		b, _ := e2.(*ast.FloatLiteral)
		return &ast.Boolean{
			Value: a.Value != b.Value,
		}, nil
	})

	// STRING != STRING
	registerInfixEvalFn(ast.STRING, "!=", ast.STRING, func(e1 ast.Expression, s string, e2 ast.Expression) (ast.Expression, error) {
		a, _ := e1.(*ast.StringLiteral)
		b, _ := e2.(*ast.StringLiteral)
		return &ast.Boolean{
			Value: a.Value != b.Value,
		}, nil
	})

	// INTEGER < INTEGER
	registerInfixEvalFn(ast.INTEGER, "<", ast.INTEGER, func(e1 ast.Expression, s string, e2 ast.Expression) (ast.Expression, error) {
		a, _ := e1.(*ast.IntegerLiteral)
		b, _ := e2.(*ast.IntegerLiteral)
		return &ast.Boolean{
			Value: a.Value < b.Value,
		}, nil
	})

	// FLOAT < FLOAT
	registerInfixEvalFn(ast.FLOAT, "<", ast.FLOAT, func(e1 ast.Expression, s string, e2 ast.Expression) (ast.Expression, error) {
		a, _ := e1.(*ast.FloatLiteral)
		b, _ := e2.(*ast.FloatLiteral)
		return &ast.Boolean{
			Value: a.Value < b.Value,
		}, nil
	})

	// INTEGER > INTEGER
	registerInfixEvalFn(ast.INTEGER, ">", ast.INTEGER, func(e1 ast.Expression, s string, e2 ast.Expression) (ast.Expression, error) {
		a, _ := e1.(*ast.IntegerLiteral)
		b, _ := e2.(*ast.IntegerLiteral)
		return &ast.Boolean{
			Value: a.Value > b.Value,
		}, nil
	})

	// FLOAT > FLOAT
	registerInfixEvalFn(ast.FLOAT, ">", ast.FLOAT, func(e1 ast.Expression, s string, e2 ast.Expression) (ast.Expression, error) {
		a, _ := e1.(*ast.FloatLiteral)
		b, _ := e2.(*ast.FloatLiteral)
		return &ast.Boolean{
			Value: a.Value > b.Value,
		}, nil
	})

	// BOOLEAN && BOOLEAN
	registerInfixEvalFn(ast.BOOLEAN, "AND", ast.BOOLEAN, func(e1 ast.Expression, s string, e2 ast.Expression) (ast.Expression, error) {
		a, _ := e1.(*ast.Boolean)
		b, _ := e2.(*ast.Boolean)
		return &ast.Boolean{
			Value: a.Value && b.Value,
		}, nil
	})

	// BOOLEAN || BOOLEAN
	registerInfixEvalFn(ast.BOOLEAN, "OR", ast.BOOLEAN, func(e1 ast.Expression, s string, e2 ast.Expression) (ast.Expression, error) {
		a, _ := e1.(*ast.Boolean)
		b, _ := e2.(*ast.Boolean)
		return &ast.Boolean{
			Value: a.Value || b.Value,
		}, nil
	})
}

type Scope struct {
	vars map[string]ast.Expression
}

func NewScope() *Scope {
	return &Scope{
		vars: map[string]ast.Expression{},
	}
}

func (s *Scope) SetVar(key string, value ast.Expression) {
	s.vars[key] = value
}

func (s *Scope) GetVar(key string) ast.Expression {
	val, ok := s.vars[key]
	if !ok {
		return nil
	}
	return val
}

func EvalExpression(expr ast.Expression, scope *Scope) (ast.Expression, error) {
	switch node := expr.(type) {
	case *ast.StringLiteral:
		return node, nil
	case *ast.FloatLiteral:
		return node, nil
	case *ast.IntegerLiteral:
		return node, nil
	case *ast.Boolean:
		return node, nil
	case *ast.Identifier:
		return scope.GetVar(node.Value), nil
	case *ast.InfixExpression:
		left, err := EvalExpression(node.Left, scope)
		if err != nil {
			return nil, err
		}

		right, err := EvalExpression(node.Right, scope)
		if err != nil {
			return nil, err
		}

		fn, ok := infixEvalFns[toFnString(left, node.Operator, right)]
		if !ok {
			return nil, errors.New("invalid operation")
		}

		result, err := fn(left, node.Operator, right)
		if err != nil {
			return nil, err
		}

		return result, nil
	}

	return nil, errors.New("invalid expression")
}

func toFnString(left ast.Expression, op string, right ast.Expression) string {
	return fmt.Sprintf("%s_%s_%s", left.Type(), op, right.Type())
}

func registerInfixEvalFn(left ast.NodeType, op string, right ast.NodeType, fn infixEvalFn) {
	infixEvalFns[fmt.Sprintf("%s_%s_%s", left, op, right)] = fn
	// Add an evaluator for lowercase operators too
	infixEvalFns[fmt.Sprintf("%s_%s_%s", left, strings.ToLower(op), right)] = fn
}
