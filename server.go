package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const jsonContentType = "application/json"

type TodoStore interface {
	GetTodo(id string) Todo
	GetTodos() Todos
	AddTodo(todo Todo)
}

type TodoServer struct {
	store TodoStore
	http.Handler
}

func NewTodoServer(store TodoStore) *TodoServer {
	s := new(TodoServer)

	s.store = store

	router := http.NewServeMux()
	router.Handle("/todos/", http.HandlerFunc(s.todoHandler))
	router.Handle("/todos", http.HandlerFunc(s.todosHandler))

	s.Handler = router

	return s
}

func (s *TodoServer) todoHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/todos/"):]

	todo := s.store.GetTodo(id)
	if todo.Id == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("content-type", jsonContentType)
	err := json.NewEncoder(w).Encode(todo)
	if err != nil {
		log.Fatal("Could not encode Todo into JSON", err)
	}
}

func (s *TodoServer) todosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("content-type", jsonContentType)
		err := json.NewEncoder(w).Encode(s.store.GetTodos())
		if err != nil {
			log.Fatal("Could not encode Todos into JSON", err)
		}
	case http.MethodPost:
		var todo Todo
		err := json.NewDecoder(r.Body).Decode(&todo)
		if err != nil {
			log.Fatal("Could not decode JSON", err)
		}

		s.store.AddTodo(todo)

		w.WriteHeader(http.StatusAccepted)
	}
}
