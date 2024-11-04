package main

import (
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request, route string) {
	w.Write([]byte("Hello from GET API route"))
}
