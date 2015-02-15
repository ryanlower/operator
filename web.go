package main

import (
	"log"
	"net/http"
)

func main() {
	operator := new(Operator)
	operator.config = new(Config)
	operator.config.load()

	log.Printf("Operator listening on port %v ...", operator.config.port)
	err := http.ListenAndServe(":"+operator.config.port, operator)
	if err != nil {
		panic(err)
	}
}
