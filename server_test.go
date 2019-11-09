package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestTodoServer(t *testing.T) {
	server := NewTodoServer(&StubTodoStore{
		[]Todo{
			{"id1", "meet friend"},
		},
	})

	t.Run("hello world", func(t *testing.T) {
		request := newGetRequest("id1")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var todo Todo
		err := json.NewDecoder(response.Body).Decode(&todo)
		if err != nil {
			t.Errorf("Decoding Todo JSON failed, %v", err)
		}

		want := Todo{"id1", "meet friend"}
		if !reflect.DeepEqual(todo, want) {
			t.Errorf("got %v want %v", todo, want)
		}
	})
}

func newGetRequest(id string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/todos/%s", id), nil)
	return req
}
