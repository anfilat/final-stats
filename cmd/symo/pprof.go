// +build pprof

// go tool pprof -http=":8080" http://127.0.0.1:7070/debug/pprof/profile?seconds=30
// go tool pprof -http=":8080" http://127.0.0.1:7070/debug/pprof/heap

package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

func init() {
	go func() {
		log.Println(http.ListenAndServe("localhost:7070", nil))
	}()
}
