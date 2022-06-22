package main

import (
	_ "embed"
	"fmt"
	"net/http"
)

//go:embed frontend/index.html
var root_page string

func RootRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, root_page)
}
