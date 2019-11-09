package main

import (
	"fmt"
	"net/http"
)

type Todo struct {
	Id   string
	Task string
}

type Server struct {
	http.Handler
}

func NewServer() *Server {
	s := new(Server)

	router := http.NewServeMux()
	router.Handle("/todos/", http.HandlerFunc(s.todosHandler))

	s.Handler = router

	return s
}

func (s *Server) todosHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello world")
}
