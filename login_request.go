package main

import (
	_ "embed"
	"fmt"
	"net/http"
)

//go:embed frontend/index.html
var login_page string

func LoginRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, login_page)
}
