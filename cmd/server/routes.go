package main

import (
	"encoding/json"
	"fmt"
	"jnafolayan/sql-db/lexer"
	"jnafolayan/sql-db/parser"
	"net/http"
)

type ErrResponse struct {
	Error string `json:"error"`
	Code  uint   `json:"code"`
}

func handleQueryText(w http.ResponseWriter, r *http.Request) {
	stmts := r.URL.Query().Get("q")

	l := lexer.New(stmts)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		resp := &ErrResponse{
			Error: err.Error(),
			Code:  http.StatusOK,
		}

		json.NewEncoder(w).Encode(resp)
		return
	}

	fmt.Println(program)
}
