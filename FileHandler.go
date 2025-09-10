package main

import "net/http"

func FileHandler() http.Handler {
	return http.FileServer(http.Dir("./static/"))
}
