package services

import (
	"fmt"

	"github.com/lloydmeta/todddo-openapi/internal/domain"
)

type TodoService interface {
	Create(newTodo *domain.NewTodo) (domain.Todo, TodoServiceError)
	Update(todo *domain.Todo) (domain.Todo, TodoServiceError)
	List() []domain.Todo
	Get(todoId *domain.TodoID) (domain.Todo, TodoServiceError)
	Delete(todoId *domain.TodoID) (bool, TodoServiceError)
}

// MkTodoService returns a default implementation of TodoService given
// a domain.TodoRepo
func MkTodoService(repo domain.TodoRepo) TodoService {
	return &todoServiceImpl{Repo: repo}
}

// todoServiceImpl encapsulates business logic around domain.Todo
//
// At the moment it mostly does simple validation on Create and Update
// cases, but maybe it can do more interesting things in the future
type todoServiceImpl struct {
	Repo domain.TodoRepo
}

func (service *todoServiceImpl) Create(newTodo *domain.NewTodo) (domain.Todo, TodoServiceError) {
	if len(newTodo.Task) == 0 {
		err := &TodoDataError{Task: newTodo.Task}
		return domain.Todo{}, err
	} else {
		return service.Repo.Create(newTodo), nil
	}
}

func (service *todoServiceImpl) Update(todo *domain.Todo) (domain.Todo, TodoServiceError) {
	if len(todo.Task) == 0 {
		err := &TodoDataError{Task: todo.Task}
		return domain.Todo{}, err
	} else {
		if updated, err := service.Repo.Update(todo); err == nil {
			return updated, nil
		} else {
			return domain.Todo{}, TodoNotFound{ID: err.Id()}
		}
	}
}

func (service *todoServiceImpl) List() []domain.Todo {
	return service.Repo.List()
}

func (service *todoServiceImpl) Get(todoId *domain.TodoID) (domain.Todo, TodoServiceError) {
	if found, err := service.Repo.Get(todoId); err == nil {
		return found, nil
	} else {
		return domain.Todo{}, TodoNotFound{ID: err.Id()}
	}
}

func (service *todoServiceImpl) Delete(todoId *domain.TodoID) (bool, TodoServiceError) {
	if result, err := service.Repo.Delete(todoId); err == nil {
		return result, nil
	} else {
		return false, TodoNotFound{ID: err.Id()}
	}
}

// <-- errors

type TodoServiceError interface {
	error
}

type TodoDataError struct {
	Task string
}

type TodoNotFound struct {
	ID domain.TodoID
}

func (err TodoDataError) Error() string {
	return fmt.Sprintf("This task was empty: [%s]", err.Task)
}

func (err TodoNotFound) Error() string {
	return fmt.Sprintf("This id does not exist: [%v]", err.ID)
}

//     errors  -->
