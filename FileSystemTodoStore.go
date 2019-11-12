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
	err := initializeTodoDBFile(file)
	if err != nil {
		return nil, fmt.Errorf("problem initialising todo db file, %v", err)
	}

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

func initializeTodoDBFile(file *os.File) error {
	file.Seek(0, 0)
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, 0)
	}

	return nil
}
