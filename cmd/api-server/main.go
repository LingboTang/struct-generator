// Command api-server serves struct-generator's functionality over HTTP.
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/LingboTang/struct-generator/api"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":8080", "address to listen on")
	flag.Parse()

	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, api.NewHandler()); err != nil {
		log.Fatal(err)
	}
}
