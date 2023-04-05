package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

const ADDR = ":8080"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/query", handleQueryText)

	err := http.ListenAndServe(ADDR, mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Fprintf(os.Stderr, "server closed: %s", err)
	} else {
		fmt.Fprintf(os.Stderr, "init error: %s", err)
	}
}
