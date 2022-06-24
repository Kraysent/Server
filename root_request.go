package main

import (
	_ "embed"
	"fmt"
	"net/http"
)

var root_page string

func RootRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, root_page)
}
