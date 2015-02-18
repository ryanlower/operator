package main

import (
	"log"
	"net/http"
)

func main() {
	operator := new(Operator)
	operator.config = new(Config)
	operator.config.load()

	log.Printf("Operator listening on port %v ...", operator.config.Port)
	err := http.ListenAndServe(":"+operator.config.Port, operator)
	if err != nil {
		panic(err)
	}
}
