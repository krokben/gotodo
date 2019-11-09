package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	server := NewServer()
	t.Run("hello world", func(t *testing.T) {
		request := newGetRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "hello world"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}

func newGetRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/todos/", nil)
	return req
}
