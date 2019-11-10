package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestTodoServerIntegrationTest(t *testing.T) {
	server := NewTodoServer(&StubTodoStore{})

	server.ServeHTTP(httptest.NewRecorder(), newPostRequest(t, "/todos", Todo{"id1", "meet friend"}))
	server.ServeHTTP(httptest.NewRecorder(), newPostRequest(t, "/todos", Todo{"id2", "buy snacks"}))

	t.Run("POST then GET todos", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetRequest("/todos"))

		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, jsonContentType)

		var todos Todos
		err := json.NewDecoder(response.Body).Decode(&todos)
		if err != nil {
			t.Errorf("Decoding Todos JSON failed, %v", err)
		}

		want := Todos{
			{"id1", "meet friend"},
			{"id2", "buy snacks"},
		}
		if !reflect.DeepEqual(todos, want) {
			t.Errorf("got %v want %v", todos, want)
		}
	})
}
