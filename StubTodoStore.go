package main

type StubTodoStore struct {
	todos Todos
}

func (s *StubTodoStore) GetTodo(id string) Todo {
	return s.GetTodos().Find(id)
}

func (s *StubTodoStore) GetTodos() Todos {
	return s.todos
}

func (s *StubTodoStore) AddTodo(todo Todo) {
	s.todos = append(s.todos, todo)
}
