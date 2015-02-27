package main

import (
	"log"
	"net/http"
)

// Operator is a ...
type Operator struct {
	config *Config
	store  Store
}

func (o *Operator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if !o.authenticated(r) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		o.create(w, r)
	} else {
		o.lookup(w, r)
	}
}

// Is this request authenticated?
// Returns true if AUTH_PASSWORD is set and provided password matches
// or if AUTH_PASSWORD is not set
// Returns false if AUTH_PASSWORD is set and password doesn't match
func (o *Operator) authenticated(r *http.Request) bool {
	log.Print(o.config.Auth.Password)

	_, password, _ := r.BasicAuth()
	if o.config.Auth.Password != "" && o.config.Auth.Password != password {
		return false
	}
	return true
}

// Create token in store, with url as value
// Returns 200 if token created successfully
// Returns 400 bad request if not
func (o *Operator) create(w http.ResponseWriter, r *http.Request) {
	token := o.parseToken(r)
	url := r.FormValue("url")

	err := o.store.Set(token, url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Lookup url by token in store
// Redirects to url if found
// Returns 404 not found if not
func (o *Operator) lookup(w http.ResponseWriter, r *http.Request) {
	token := o.parseToken(r)

	url, err := o.store.Get(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Printf("Connecting %v to %v", token, url)
	http.Redirect(w, r, url, 301)
}

func (o *Operator) parseToken(r *http.Request) string {
	return r.URL.Path[1:] // Strip leading slash
}
