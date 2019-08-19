package controllers

import (
	"fmt"
	"github.com/lloydmeta/todddo-openapi/internal/api/models"
	"github.com/lloydmeta/todddo-openapi/internal/domain"
	"github.com/lloydmeta/todddo-openapi/internal/domain/services"
	"net/http"
)

type TodoController interface {
	Create(newTodo *models.TodoData) (models.Todo, models.ApiError)
	Get(id *domain.TodoID) (models.Todo, models.ApiError)
	Delete(id *domain.TodoID) (models.Success, models.ApiError)
	List() []models.Todo
	Update(todo *models.Todo) (models.Todo, models.ApiError)
}

// MkTodosController returns a TodoController when given a services.TodoService
func MkTodosController(service services.TodoService) TodoController {
	return &TodosControllerImpl{service: service}
}

type TodosControllerImpl struct {
	service services.TodoService
}

func (t *TodosControllerImpl) Create(newTodo *models.TodoData) (models.Todo, models.ApiError) {
	domainTodo := domain.NewTodo{
		Task: newTodo.Task,
	}
	if persisted, err := t.service.Create(&domainTodo); err == nil {
		return toApiTodo(&persisted), nil
	} else {
		return models.Todo{}, TodosControllerError{
			httpStatusCode: http.StatusBadRequest,
			message:        err.Error(),
		}
	}

}

func (t *TodosControllerImpl) Get(id *domain.TodoID) (models.Todo, models.ApiError) {
	if found, err := t.service.Get(id); err == nil {
		return toApiTodo(&found), nil
	} else {
		return models.Todo{}, TodosControllerError{
			httpStatusCode: http.StatusNotFound,
			message:        err.Error(),
		}
	}
}

func (t *TodosControllerImpl) Delete(id *domain.TodoID) (models.Success, models.ApiError) {
	if _, err := t.service.Delete(id); err == nil {
		return models.Success{Message: fmt.Sprintf("Successfully deleted Todo with id [%v]", *id)}, nil
	} else {
		return models.Success{}, TodosControllerError{
			httpStatusCode: http.StatusNotFound,
			message:        err.Error(),
		}
	}
}

func (t *TodosControllerImpl) List() []models.Todo {
	domainTodos := t.service.List()
	apiTodos := make([]models.Todo, len(domainTodos))
	for i, domainTodo := range domainTodos {
		apiTodos[i] = toApiTodo(&domainTodo)
	}
	return apiTodos
}

func (t *TodosControllerImpl) Update(todo *models.Todo) (models.Todo, models.ApiError) {
	domainTodo := toDomainTodo(todo)
	if _, err := t.service.Update(&domainTodo); err == nil {
		return *todo, nil
	} else {
		switch err.(type) {
		case services.TodoNotFound:
			return *todo, TodosControllerError{
				httpStatusCode: http.StatusNotFound,
				message:        err.Error(),
			}
		default:
			return *todo, TodosControllerError{
				httpStatusCode: http.StatusBadRequest,
				message:        err.Error(),
			}
		}
	}
}

func toApiTodo(domainTodo *domain.Todo) models.Todo {
	return models.Todo{
		ID:   domainTodo.ID,
		Task: domainTodo.Task,
	}
}
func toDomainTodo(apiTodo *models.Todo) domain.Todo {
	return domain.Todo{
		ID:   apiTodo.ID,
		Task: apiTodo.Task,
	}
}

type TodosControllerError struct {
	httpStatusCode int
	message        string
}

func (t TodosControllerError) Error() string {
	return t.message
}

func (t TodosControllerError) AsModel() models.Error {
	return models.Error{Message: t.message}
}

func (t TodosControllerError) HttpStatusCode() int {
	return t.httpStatusCode
}
