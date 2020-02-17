package main

import (
	"net/http"
	_ "net/http/pprof"
)

func startpprof() {
	go func() {
		http.ListenAndServe("localhost:8080", nil)
	}()
}
