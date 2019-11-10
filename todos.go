package main

import (
	"encoding/json"
	"fmt"
	"io"
)

func NewTodos(rdr io.Reader) ([]Todo, error) {
	var todos []Todo
	err := json.NewDecoder(rdr).Decode(&todos)
	if err != nil {
		err = fmt.Errorf("problem parsing league, %v", err)
	}

	return todos, err
}
