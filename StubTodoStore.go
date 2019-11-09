package main

type StubTodoStore struct {
	todos []Todo
}

func (s *StubTodoStore) GetTodo(id string) Todo {
	var result Todo

	for _, todo := range s.todos {
		if todo.Id == id {
			result = todo
		}
	}

	return result
}

func (s *StubTodoStore) GetTodos() []Todo {
	return s.todos
}
