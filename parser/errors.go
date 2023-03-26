package parser

import (
	"errors"
	"fmt"
	"jnafolayan/sql-db/token"
)

func expectedTokenError(expected token.TokenType) error {
	return fmt.Errorf("expected %s", expected)
}

var ErrEmptyColumnsList = errors.New("must specify a column name")
