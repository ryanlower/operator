package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup() (*Operator, *httptest.ResponseRecorder) {
	operator := new(Operator)
	operator.connect()
	recorder := httptest.NewRecorder()

	return operator, recorder
}

func TestCreationNewToken(t *testing.T) {
	token := "ryan"
	url := "http://ryanlower.com"

	operator, w := setup()
	operator.connection.Do("DEL", token) // Ensure token doesn't exist

	req, _ := http.NewRequest("POST", "/"+token+"?url="+url, nil)
	operator.create(w, req)

	assert.Equal(t, w.Code, http.StatusOK)
}

func TestCreationExistingToken(t *testing.T) {
	token := "ryan"
	url := "http://ryanlower.com"

	operator, w := setup()
	operator.connection.Do("SET", token, url) // Ensure token exists

	req, _ := http.NewRequest("POST", "/"+token+"?url="+url, nil)
	operator.create(w, req)

	assert.Equal(t, w.Code, http.StatusBadRequest)
	assert.Equal(t, w.Body.String(), "Token already used\n")
}

func TestLookup(t *testing.T) {
	token := "good_token"

	operator, w := setup()
	operator.connection.Do("SET", token, "http://ryanlower.com")

	req, _ := http.NewRequest("GET", "/"+token, nil)
	operator.lookup(w, req)

	assert.Equal(t, w.Code, http.StatusMovedPermanently)
	assert.Equal(t, w.Header().Get("Location"), "http://ryanlower.com")
}

func TestLookupBadToken(t *testing.T) {
	token := "bad_token"

	operator, w := setup()
	operator.connection.Do("DEL", token)

	req, _ := http.NewRequest("GET", "/"+token, nil)
	operator.lookup(w, req)

	assert.Equal(t, w.Code, http.StatusNotFound)
	assert.Equal(t, w.Body.String(), "Token not found\n")
}
