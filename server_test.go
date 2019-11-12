package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestTodoServer(t *testing.T) {
	server := NewTodoServer(&StubTodoStore{
		Todos{
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

		got, _ := NewTodos(response.Body)

		want := Todos{
			{"id1", "meet friend"},
			{"id2", "buy snacks"},
		}

		assertDeepEqual(t, got, want)
	})

	t.Run("POST todo", func(t *testing.T) {
		request := newPostRequest(t, "/todos", Todo{"id3", "find keys"})
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)
	})
}

func TestFileSystemTodoStore(t *testing.T) {
	t.Run("GET todos", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"id": "id1", "task": "meet friend"},
			{"id": "id2", "task": "buy snacks"}]`)
		defer cleanDatabase()

		store := FileSystemTodoStore{database}

		got := store.GetTodos()
		want := Todos{
			{"id1", "meet friend"},
			{"id2", "buy snacks"},
		}

		assertDeepEqual(t, got, want)

		// read again
		gotAgain := store.GetTodos()
		assertDeepEqual(t, gotAgain, want)
	})

	t.Run("GET todo", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"id": "id1", "task": "meet friend"}]`)
		defer cleanDatabase()

		store := FileSystemTodoStore{database}

		got := store.GetTodo("id1")
		want := Todo{"id1", "meet friend"}

		assertDeepEqual(t, got, want)
	})

	t.Run("POST then GET todo", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"id": "id1", "task": "meet friend"}]`)
		defer cleanDatabase()

		store := FileSystemTodoStore{database}

		store.AddTodo(Todo{"id7", "go home"})

		got := store.GetTodos()
		want := Todos{
			{"id1", "meet friend"},
			{"id7", "go home"},
		}

		assertDeepEqual(t, got, want)
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

func assertDeepEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func createTempFile(t *testing.T, initialData string) (io.ReadWriteSeeker, func()) {
	t.Helper()

	tmpfile, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpfile.Write([]byte(initialData))

	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFile
}
