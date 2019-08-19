package controllers

import (
	apiModels "github.com/lloydmeta/todddo-openapi/internal/api/models"
	"github.com/lloydmeta/todddo-openapi/internal/domain"
	"github.com/lloydmeta/todddo-openapi/internal/domain/services"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCreateOk(t *testing.T) {
	mockService := mockTodoService{}
	mockService.create = func(newTodo *domain.NewTodo) (domain.Todo, services.TodoServiceError) {
		return domain.Todo{
			ID:   domain.TodoID(123),
			Task: newTodo.Task,
		}, nil
	}
	controller := MkTodosController(&mockService)
	newTodo := apiModels.TodoData{Task: "let's get cracking"}
	r, _ := controller.Create(&newTodo)
	assert.Equal(t, 1, mockService.createCalled)
	assert.Equal(t, newTodo.Task, r.Task)
}

func TestCreateInvalidData(t *testing.T) {
	mockService := mockTodoService{}
	mockService.create = func(newTodo *domain.NewTodo) (domain.Todo, services.TodoServiceError) {
		return domain.Todo{}, services.TodoDataError{Task: newTodo.Task}
	}
	controller := MkTodosController(&mockService)
	newTodo := apiModels.TodoData{Task: "let's get cracking"}
	_, err := controller.Create(&newTodo)
	if err != nil {
		assert.Equal(t, 1, mockService.createCalled)
		assert.Equal(t, http.StatusBadRequest, err.HttpStatusCode())
	} else {
		assert.Fail(t, "Expected an error")
	}
}

func TestGetOk(t *testing.T) {
	mockService := mockTodoService{}
	todoId := domain.TodoID(1234)
	domainModel := domain.Todo{
		ID:   todoId,
		Task: "lol",
	}
	mockService.get = func(todoId *domain.TodoID) (todo domain.Todo, serviceError services.TodoServiceError) {
		return domainModel, nil
	}
	controller := MkTodosController(&mockService)
	found, _ := controller.Get(&todoId)
	assert.Equal(t, 1, mockService.getCalled)
	assert.Equal(t, domainModel.Task, found.Task)
	assert.Equal(t, domainModel.ID, found.ID)
}

func TestGetNotFound(t *testing.T) {
	mockService := mockTodoService{}
	mockService.get = func(todoId *domain.TodoID) (todo domain.Todo, serviceError services.TodoServiceError) {
		return domain.Todo{}, services.TodoNotFound{ID: *todoId}
	}
	controller := MkTodosController(&mockService)
	todoId := domain.TodoID(123)
	_, err := controller.Get(&todoId)
	if err != nil {
		assert.Equal(t, 1, mockService.getCalled)
		assert.Equal(t, http.StatusNotFound, err.HttpStatusCode())
	} else {
		assert.Fail(t, "Expected an error")
	}
}

func TestDeleteOk(t *testing.T) {
	mockService := mockTodoService{}
	todoId := domain.TodoID(1234)
	mockService.delete = func(todoId *domain.TodoID) (b bool, serviceError services.TodoServiceError) {
		return true, nil
	}
	controller := MkTodosController(&mockService)
	_, err := controller.Delete(&todoId)
	assert.Equal(t, 1, mockService.deleteCalled)
	assert.Nil(t, err)
}

func TestDeleteFound(t *testing.T) {
	mockService := mockTodoService{}
	todoId := domain.TodoID(1234)
	mockService.delete = func(todoId *domain.TodoID) (b bool, serviceError services.TodoServiceError) {
		return false, services.TodoNotFound{ID: *todoId}
	}
	controller := MkTodosController(&mockService)
	_, err := controller.Delete(&todoId)
	if err != nil {
		assert.Equal(t, 1, mockService.deleteCalled)
		assert.Equal(t, http.StatusNotFound, err.HttpStatusCode())
	} else {
		assert.Fail(t, "Expected an error")
	}
}

func TestList(t *testing.T) {
	mockService := mockTodoService{}
	todoId := domain.TodoID(1234)
	domainModel := domain.Todo{
		ID:   todoId,
		Task: "lol",
	}
	domainModels := []domain.Todo{domainModel}
	mockService.list = func() []domain.Todo {
		return domainModels
	}
	controller := MkTodosController(&mockService)
	results := controller.List()
	assert.Equal(t, 1, mockService.listCalled)
	expected := make([]apiModels.Todo, len(domainModels))
	for i, v := range domainModels {
		expected[i] = toApiTodo(&v)
	}
	assert.Equal(t, expected, results)
}

func TestUpdateOk(t *testing.T) {
	mockService := mockTodoService{}
	todoId := domain.TodoID(1234)
	apiModel := apiModels.Todo{
		ID:   todoId,
		Task: "lol",
	}
	mockService.update = func(todo *domain.Todo) (todo2 domain.Todo, serviceError services.TodoServiceError) {
		return *todo, nil
	}
	controller := MkTodosController(&mockService)
	_, err := controller.Update(&apiModel)
	assert.Equal(t, 1, mockService.updateCalled)
	assert.Nil(t, err)
}

func TestUpdateNotfound(t *testing.T) {
	mockService := mockTodoService{}
	todoId := domain.TodoID(1234)
	apiModel := apiModels.Todo{
		ID:   todoId,
		Task: "lol",
	}
	mockService.update = func(todo *domain.Todo) (todo2 domain.Todo, serviceError services.TodoServiceError) {
		return *todo, services.TodoNotFound{ID: todo.ID}
	}
	controller := MkTodosController(&mockService)
	_, err := controller.Update(&apiModel)
	if err != nil {
		assert.Equal(t, http.StatusNotFound, err.HttpStatusCode())
	} else {
		assert.Fail(t, "Expected an error")
	}
}

func TestUpdateInvaliddata(t *testing.T) {
	mockService := mockTodoService{}
	todoId := domain.TodoID(1234)
	apiModel := apiModels.Todo{
		ID:   todoId,
		Task: "lol",
	}
	mockService.update = func(todo *domain.Todo) (todo2 domain.Todo, serviceError services.TodoServiceError) {
		return *todo, services.TodoDataError{Task: todo.Task}
	}
	controller := MkTodosController(&mockService)
	_, err := controller.Update(&apiModel)
	if err != nil {
		assert.Equal(t, http.StatusBadRequest, err.HttpStatusCode())
	} else {
		assert.Fail(t, "Expected an error")
	}
}

// Mocks

type mockTodoService struct {
	create       func(newTodo *domain.NewTodo) (domain.Todo, services.TodoServiceError)
	createCalled int
	update       func(todo *domain.Todo) (domain.Todo, services.TodoServiceError)
	updateCalled int
	list         func() []domain.Todo
	listCalled   int
	get          func(todoId *domain.TodoID) (domain.Todo, services.TodoServiceError)
	getCalled    int
	delete       func(todoId *domain.TodoID) (bool, services.TodoServiceError)
	deleteCalled int
}

func (m *mockTodoService) Create(newTodo *domain.NewTodo) (domain.Todo, services.TodoServiceError) {
	defer func() { m.createCalled++ }()
	return m.create(newTodo)
}

func (m *mockTodoService) Update(todo *domain.Todo) (domain.Todo, services.TodoServiceError) {
	defer func() { m.updateCalled++ }()
	return m.update(todo)
}

func (m *mockTodoService) List() []domain.Todo {
	defer func() { m.listCalled++ }()
	return m.list()
}

func (m *mockTodoService) Get(todoId *domain.TodoID) (domain.Todo, services.TodoServiceError) {
	defer func() { m.getCalled++ }()
	return m.get(todoId)
}

func (m *mockTodoService) Delete(todoId *domain.TodoID) (bool, services.TodoServiceError) {
	defer func() { m.deleteCalled++ }()
	return m.delete(todoId)
}
