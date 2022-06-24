package main

import (
	_ "embed"
	"fmt"
	"net/http"
)

var login_page string

func LoginRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, login_page)
}
