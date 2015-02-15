package main

import (
	// "log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup() (*Operator, *httptest.ResponseRecorder) {
	operator := new(Operator)
	operator.config = new(Config)
	operator.config.redis.port = "6379"
	operator.connect()
	recorder := httptest.NewRecorder()

	return operator, recorder
}

// Add link to redis
func addLink(o *Operator, token, url string) {
	o.connection.Do("SET", token, url)
}

// Remove link in redis
func removeLink(o *Operator, token string) {
	o.connection.Do("DEL", token)
}

func TestCreationNewToken(t *testing.T) {
	token := "ryan"
	url := "http://ryanlower.com"

	operator, w := setup()
	removeLink(operator, token)

	req, _ := http.NewRequest("POST", "/"+token+"?url="+url, nil)
	operator.create(w, req)

	assert.Equal(t, w.Code, http.StatusOK)
}

func TestCreationExistingToken(t *testing.T) {
	token := "ryan"
	url := "http://ryanlower.com"

	operator, w := setup()
	addLink(operator, token, url)

	req, _ := http.NewRequest("POST", "/"+token+"?url="+url, nil)
	operator.create(w, req)

	assert.Equal(t, w.Code, http.StatusBadRequest)
	assert.Equal(t, w.Body.String(), "Token already used\n")
}

func TestLookup(t *testing.T) {
	token := "good_token"

	operator, w := setup()
	addLink(operator, token, "http://ryanlower.com")

	req, _ := http.NewRequest("GET", "/"+token, nil)
	operator.lookup(w, req)

	assert.Equal(t, w.Code, http.StatusMovedPermanently)
	assert.Equal(t, w.Header().Get("Location"), "http://ryanlower.com")
}

func TestLookupBadToken(t *testing.T) {
	token := "bad_token"

	operator, w := setup()
	removeLink(operator, token)

	req, _ := http.NewRequest("GET", "/"+token, nil)
	operator.lookup(w, req)

	assert.Equal(t, w.Code, http.StatusNotFound)
	assert.Equal(t, w.Body.String(), "Token not found\n")
}

func TestCreationAuthenticationBadAuth(t *testing.T) {
	operator, w := setup()
	operator.config.auth.password = "password"
	removeLink(operator, "token")

	req, _ := http.NewRequest("POST", "/token?url=url", nil)
	operator.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreationAuthenticationGoodAuth(t *testing.T) {
	operator, w := setup()
	operator.config.auth.password = "password"
	removeLink(operator, "token")

	req, _ := http.NewRequest("POST", "/token?url=url", nil)
	req.SetBasicAuth("", "password")
	operator.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
