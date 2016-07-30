// +build !appengine

package main

import (
	"github.com/fatlotus/ironclad"
	"net/http"
)

func main() {
	panic(http.ListenAndServe(":9000", ironclad.New()))
}
