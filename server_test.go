package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestTodoServer(t *testing.T) {
	server := NewTodoServer(&StubTodoStore{
		[]Todo{
			{"id1", "meet friend"},
			{"id2", "buy snacks"},
		},
	})

	t.Run("GET todo", func(t *testing.T) {
		request := newGetRequest("/todos/id1")
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

	t.Run("GET todos", func(t *testing.T) {
		request := newGetRequest("/todos")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, jsonContentType)

		var todos []Todo
		err := json.NewDecoder(response.Body).Decode(&todos)
		if err != nil {
			t.Errorf("Decoding Todos JSON failed, %v", err)
		}

		want := []Todo{
			{"id1", "meet friend"},
			{"id2", "buy snacks"},
		}
		if !reflect.DeepEqual(todos, want) {
			t.Errorf("got %v want %v", todos, want)
		}
	})

	t.Run("POST todo", func(t *testing.T) {
		request := newPostRequest(t, "/todos", Todo{"id3", "find keys"})
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)
	})
}

func newGetRequest(endpoint string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
	return req
}

func newPostRequest(t *testing.T, endpoint string, data Todo) *http.Request {
	t.Helper()

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Could not marshal data, %v", err)
	}

	req, _ := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(jsonData))
	return req
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got status %d want %d", got, want)
	}
}

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}
