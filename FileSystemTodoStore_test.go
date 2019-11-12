package main

import "testing"

func TestFileSystemTodoStore(t *testing.T) {
	t.Run("GET todos", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"id": "id1", "task": "meet friend"},
			{"id": "id2", "task": "buy snacks"}]`)
		defer cleanDatabase()

		store := NewFileSystemTodoStore(database)

		got := store.GetTodos()
		want := Todos{
			{"id1", "meet friend"},
			{"id2", "buy snacks"},
		}

		assertDeepEqual(t, got, want)

		// read again
		gotAgain := store.GetTodos()
		assertDeepEqual(t, gotAgain, want)
	})

	t.Run("GET todo", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"id": "id1", "task": "meet friend"}]`)
		defer cleanDatabase()

		store := NewFileSystemTodoStore(database)

		got := store.GetTodo("id1")
		want := Todo{"id1", "meet friend"}

		assertDeepEqual(t, got, want)
	})

	t.Run("POST then GET todo", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"id": "id1", "task": "meet friend"}]`)
		defer cleanDatabase()

		store := NewFileSystemTodoStore(database)

		store.AddTodo(Todo{"id7", "go home"})

		got := store.GetTodos()
		want := Todos{
			{"id1", "meet friend"},
			{"id7", "go home"},
		}

		assertDeepEqual(t, got, want)
	})
}
