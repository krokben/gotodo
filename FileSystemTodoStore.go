package main

import (
	"io"
)

type FileSystemTodoStore struct {
	database io.ReadSeeker
}

func (f *FileSystemTodoStore) GetTodos() []Todo {
	f.database.Seek(0, 0)
	todos, _ := NewTodos(f.database)
	return todos
}
