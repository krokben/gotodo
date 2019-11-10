package main

import (
	"io"
)

type FileSystemTodoStore struct {
	database io.ReadSeeker
}

func (f *FileSystemTodoStore) GetTodo(id string) Todo {
	var result Todo

	todos, _ := NewTodos(f.database)
	for _, todo := range todos {
		if todo.Id == id {
			result = todo
		}
	}

	return result
}

func (f *FileSystemTodoStore) GetTodos() []Todo {
	f.database.Seek(0, 0)
	todos, _ := NewTodos(f.database)
	return todos
}
