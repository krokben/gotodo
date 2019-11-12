package main

import (
	"encoding/json"
	"io"
	"log"
)

type FileSystemTodoStore struct {
	database io.ReadWriteSeeker
	todos    Todos
}

func NewFileSystemTodoStore(database io.ReadWriteSeeker) *FileSystemTodoStore {
	database.Seek(0, 0)
	todos, _ := NewTodos(database)
	return &FileSystemTodoStore{
		database: database,
		todos:    todos,
	}
}

func (f *FileSystemTodoStore) GetTodo(id string) Todo {
	return f.GetTodos().Find(id)
}

func (f *FileSystemTodoStore) GetTodos() Todos {
	return f.todos
}

func (f *FileSystemTodoStore) AddTodo(todo Todo) {
	f.todos = append(f.todos, todo)

	f.database.Seek(0, 0)
	err := json.NewEncoder(f.database).Encode(f.todos)
	if err != nil {
		log.Fatal("Could not encode into JSON", err)
	}
}
