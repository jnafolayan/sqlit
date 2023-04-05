package main

import (
	"jnafolayan/sql-db/repl"
	"os"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
