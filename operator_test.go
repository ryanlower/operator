package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStore struct {
	Store
	data map[string]string
}

func (s *TestStore) connect() {
	s.data = map[string]string{}
}

func (s TestStore) Get(token string) (string, error) {
	url, ok := s.data[token]
	if !ok {
		return url, errors.New("Token not found")
	}
	return url, nil
}

func (s TestStore) Set(token, url string) error {
	_, ok := s.data[token]
	if ok {
		return errors.New("Token already used")
	}
	s.data[token] = url
	return nil
}

func setup() (*Operator, *httptest.ResponseRecorder) {
	conf := new(Config)

	store := new(TestStore) // A new TestStore will always be empty
	store.connect()

	operator := &Operator{config: conf, store: store}

	recorder := httptest.NewRecorder()

	return operator, recorder
}

// Add link to store
func addLink(o *Operator, token, url string) {
	o.store.Set(token, url)
}

func TestCreationNewToken(t *testing.T) {
	token := "ryan"
	url := "http://ryanlower.com"

	operator, w := setup()

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

	req, _ := http.NewRequest("GET", "/"+token, nil)
	operator.lookup(w, req)

	assert.Equal(t, w.Code, http.StatusNotFound)
	assert.Equal(t, w.Body.String(), "Token not found\n")
}

func TestCreationAuthenticationBadAuth(t *testing.T) {
	operator, w := setup()
	operator.config.Auth.Password = "password"

	req, _ := http.NewRequest("POST", "/token?url=url", nil)
	operator.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreationAuthenticationGoodAuth(t *testing.T) {
	operator, w := setup()
	operator.config.Auth.Password = "password"

	req, _ := http.NewRequest("POST", "/token?url=url", nil)
	req.SetBasicAuth("", "password")
	operator.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
