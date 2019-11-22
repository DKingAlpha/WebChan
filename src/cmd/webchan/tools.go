package main

import (
	"net/http"
	"strings"
)

func toolHandler(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(r.URL.Path, "/")[1:]
	//q := GetUrlArgs(r.URL.RawQuery)
	if len(p) < 2 {

	}
}