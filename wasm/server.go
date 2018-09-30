// +build disabled
// use: go run server.go

package main

import (
	"log"
	"net/http"
)

func main() {
	const wasm = "/cat-o-licious.wasm"
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc(wasm, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/wasm")
		http.ServeFile(w, r, "."+wasm)
	})
	log.Fatal(http.ListenAndServe(":8000", nil))
}
