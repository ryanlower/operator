package main

import (
	"log"
	"net/http"
	"os"
)

func envOrDefault(key string, defaultValue string) string {
	if env := os.Getenv(key); env != "" {
		return env
	}
	return defaultValue
}

func main() {
	operator := new(Operator)

	port := envOrDefault("PORT", "3000")
	log.Printf("Operator listening on port %v ...", port)
	err := http.ListenAndServe(":"+port, operator)
	if err != nil {
		panic(err)
	}
}
