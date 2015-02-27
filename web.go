package main

import (
	"log"
	"net/http"
)

func main() {
	conf := new(Config)
	conf.load()

	store := &RedisStore{config: conf}
	operator := &Operator{config: conf, store: store}

	log.Printf("Operator listening on port %v ...", conf.Port)
	err := http.ListenAndServe(":"+conf.Port, operator)
	if err != nil {
		panic(err)
	}
}
