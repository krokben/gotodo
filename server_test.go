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

		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, jsonContentType)

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

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
	}
}

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}
