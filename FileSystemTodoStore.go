package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type FileSystemTodoStore struct {
	database io.Writer
	todos    Todos
}

func NewFileSystemTodoStore(file *os.File) (*FileSystemTodoStore, error) {
	file.Seek(0, 0)
	todos, err := NewTodos(file)
	if err != nil {
		return nil, fmt.Errorf("problem loading todo store from file %s, %v", file.Name(), err)
	}

	return &FileSystemTodoStore{
		database: &tape{file},
		todos:    todos,
	}, nil
}

func (f *FileSystemTodoStore) GetTodo(id string) Todo {
	return f.GetTodos().Find(id)
}

func (f *FileSystemTodoStore) GetTodos() Todos {
	return f.todos
}

func (f *FileSystemTodoStore) AddTodo(todo Todo) {
	f.todos = append(f.todos, todo)

	err := json.NewEncoder(f.database).Encode(f.todos)
	if err != nil {
		log.Fatal("Could not encode into JSON", err)
	}
}
