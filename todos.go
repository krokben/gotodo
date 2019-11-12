package main

import (
	"encoding/json"
	"fmt"
	"io"
)

type Todo struct {
	Id   string `json:"id"`
	Task string `json:"task"`
}

type Todos []Todo

func NewTodos(rdr io.Reader) (Todos, error) {
	var todos Todos
	err := json.NewDecoder(rdr).Decode(&todos)
	if err != nil {
		err = fmt.Errorf("problem parsing todos, %v", err)
	}

	return todos, err
}

func (t Todos) Find(id string) Todo {
	var result Todo

	for _, todo := range t {
		if todo.Id == id {
			return todo
		}
	}

	return result
}
