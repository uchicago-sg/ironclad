// +build !appengine

package main

import (
	"github.com/fatlotus/ironclad"
	"log"
	"net/http"
)

func main() {
	log.Print("starting server on http://0.0.0.0:8080")
	log.Fatal(http.ListenAndServe(":8080", ironclad.New()))
}
