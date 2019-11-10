package main

import (
	"encoding/json"
	"io"
	"log"
)

type FileSystemTodoStore struct {
	database io.ReadWriteSeeker
}

func (f *FileSystemTodoStore) GetTodo(id string) Todo {
	return f.GetTodos().Find(id)
}

func (f *FileSystemTodoStore) GetTodos() Todos {
	f.database.Seek(0, 0)
	todos, _ := NewTodos(f.database)
	return todos
}

func (f *FileSystemTodoStore) AddTodo(todo Todo) {
	todos := f.GetTodos()
	todos = append(todos, todo)

	f.database.Seek(0, 0)
	err := json.NewEncoder(f.database).Encode(todos)
	if err != nil {
		log.Fatal("Could not encode into JSON", err)
	}
}
