package routing

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/lloydmeta/todddo-openapi/internal/api/models"
	"github.com/lloydmeta/todddo-openapi/internal/domain"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupRouter() (*gin.Engine, *mockTodoController) {
	engine := gin.Default()
	mockController := mockTodoController{}
	handler := TodosRoutesHandler{Controller: &mockController}
	handler.RegisterRoutes(engine)

	return engine, &mockController
}

func performRequest(r http.Handler, method, url string, body interface{}) *httptest.ResponseRecorder {
	var bodyToSend io.Reader
	if body != nil {
		asBytes, _ := json.Marshal(body)
		bodyToSend = bytes.NewBuffer(asBytes)
	}
	req, _ := http.NewRequest(method, url, bodyToSend)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestPostTasksOk(t *testing.T) {
	router, mockController := setupRouter()
	mockController.create = func(newTodo *models.TodoData) (todo models.Todo, apiError models.ApiError) {
		return models.Todo{
			ID:   domain.TodoID(123),
			Task: newTodo.Task,
		}, nil
	}
	newTodo := models.TodoData{Task: "do something"}
	resp := performRequest(router, http.MethodPost, "/tasks", newTodo)
	assert.Equal(t, http.StatusCreated, resp.Code)
	var respTask models.Todo
	if err := json.Unmarshal(resp.Body.Bytes(), &respTask); err != nil {
		assert.Fail(t, err.Error())
	} else {
		assert.Equal(t, newTodo.Task, respTask.Task)
		assert.Equal(t, 1, mockController.createCalled)
	}
}

func TestPostTasksInvalid(t *testing.T) {
	router, mockController := setupRouter()
	mockController.create = func(newTodo *models.TodoData) (todo models.Todo, apiError models.ApiError) {
		return models.Todo{}, mockApiError{
			code:    http.StatusBadRequest,
			message: "invalid data",
		}
	}
	newTodo := models.TodoData{Task: ""}
	resp := performRequest(router, http.MethodPost, "/tasks", newTodo)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	// the binding is `required`, so we fail at unmarshalling, which is nice
	assert.Equal(t, 0, mockController.createCalled)
}

func TestGetOk(t *testing.T) {
	router, mockController := setupRouter()
	expected := models.Todo{
		Task: "something",
	}
	mockController.get = func(id *domain.TodoID) (todo models.Todo, apiError models.ApiError) {
		expected.ID = *id
		return expected, nil
	}
	resp := performRequest(router, http.MethodGet, "/tasks/1", nil)
	assert.Equal(t, http.StatusOK, resp.Code)
	var respTask models.Todo
	if err := json.Unmarshal(resp.Body.Bytes(), &respTask); err != nil {
		assert.Fail(t, err.Error())
	} else {
		assert.Equal(t, expected.Task, respTask.Task)
		assert.Equal(t, 1, mockController.getCalled)
	}
}

func TestGetInvalidId(t *testing.T) {
	router, mockController := setupRouter()
	resp := performRequest(router, http.MethodGet, "/tasks/bababoo", nil)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, 0, mockController.getCalled)
}

func TestGetNotFound(t *testing.T) {
	router, mockController := setupRouter()
	mockController.get = func(id *domain.TodoID) (todo models.Todo, apiError models.ApiError) {
		return models.Todo{}, mockApiError{
			code:    http.StatusNotFound,
			message: "nope",
		}
	}
	resp := performRequest(router, http.MethodGet, "/tasks/1", nil)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Equal(t, 1, mockController.getCalled)
}

func TestDeleteOk(t *testing.T) {
	router, mockController := setupRouter()
	mockController.delete = func(id *domain.TodoID) (success models.Success, apiError models.ApiError) {
		return models.Success{Message: "oooh yeea"}, nil
	}
	resp := performRequest(router, http.MethodDelete, "/tasks/1", nil)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, 1, mockController.deleteCalled)
}

func TestDeleteInvalidId(t *testing.T) {
	router, mockController := setupRouter()
	resp := performRequest(router, http.MethodDelete, "/tasks/bababoo", nil)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, 0, mockController.deleteCalled)
}

func TestDeleteNotFound(t *testing.T) {
	router, mockController := setupRouter()
	mockController.delete = func(id *domain.TodoID) (success models.Success, apiError models.ApiError) {
		return models.Success{}, mockApiError{
			code:    http.StatusNotFound,
			message: "nope",
		}
	}
	resp := performRequest(router, http.MethodDelete, "/tasks/1", nil)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Equal(t, 1, mockController.deleteCalled)
}

func TestList(t *testing.T) {
	router, mockController := setupRouter()
	expected := []models.Todo{
		{
			ID:   domain.TodoID(123),
			Task: "mockity",
		},
	}
	mockController.list = func() []models.Todo {
		return expected
	}
	resp := performRequest(router, http.MethodGet, "/tasks", nil)
	var respTasks []models.Todo
	if err := json.Unmarshal(resp.Body.Bytes(), &respTasks); err != nil {
		assert.Fail(t, err.Error())
	} else {
		assert.Equal(t, expected, respTasks)
		assert.Equal(t, 1, mockController.listCalled)
	}
}

func TestUpdateTasksOk(t *testing.T) {
	router, mockController := setupRouter()
	mockController.update = func(todo *models.Todo) (todo2 models.Todo, apiError models.ApiError) {
		return *todo, nil
	}
	update := models.Todo{ID: domain.TodoID(123), Task: "do something"}
	resp := performRequest(router, http.MethodPut, "/tasks/1", update)
	assert.Equal(t, http.StatusOK, resp.Code)
	var respTask models.Todo
	if err := json.Unmarshal(resp.Body.Bytes(), &respTask); err != nil {
		assert.Fail(t, err.Error())
	} else {
		assert.Equal(t, update.Task, respTask.Task)
		assert.Equal(t, 1, mockController.updateCalled)
	}
}

func TestUpdateInvalidId(t *testing.T) {
	router, mockController := setupRouter()
	resp := performRequest(router, http.MethodPut, "/tasks/bababoo", nil)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, 0, mockController.updateCalled)
}

func TestUpdateTasksFailure(t *testing.T) {
	router, mockController := setupRouter()
	mockController.update = func(todo *models.Todo) (todo2 models.Todo, apiError models.ApiError) {
		return *todo, mockApiError{
			code:    http.StatusNotFound,
			message: "nope",
		}
	}
	update := models.Todo{ID: domain.TodoID(123), Task: "do something"}
	resp := performRequest(router, http.MethodPut, "/tasks/1", update)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Equal(t, 1, mockController.updateCalled)
}

// Mocks

type mockTodoController struct {
	create       func(newTodo *models.TodoData) (models.Todo, models.ApiError)
	createCalled int
	update       func(todo *models.Todo) (models.Todo, models.ApiError)
	updateCalled int
	list         func() []models.Todo
	listCalled   int
	get          func(id *domain.TodoID) (models.Todo, models.ApiError)
	getCalled    int
	delete       func(id *domain.TodoID) (models.Success, models.ApiError)
	deleteCalled int
}

func (m *mockTodoController) Create(newTodo *models.TodoData) (models.Todo, models.ApiError) {
	defer func() { m.createCalled++ }()
	return m.create(newTodo)
}

func (m *mockTodoController) Get(id *domain.TodoID) (models.Todo, models.ApiError) {
	defer func() { m.getCalled++ }()
	return m.get(id)
}

func (m *mockTodoController) Delete(id *domain.TodoID) (models.Success, models.ApiError) {
	defer func() { m.deleteCalled++ }()
	return m.delete(id)
}

func (m *mockTodoController) List() []models.Todo {
	defer func() { m.listCalled++ }()
	return m.list()
}

func (m *mockTodoController) Update(todo *models.Todo) (models.Todo, models.ApiError) {
	defer func() { m.updateCalled++ }()
	return m.update(todo)
}

type mockApiError struct {
	code    int
	message string
}

func (m mockApiError) Error() string {
	return m.message
}

func (m mockApiError) AsModel() models.Error {
	return models.Error{Message: m.message}
}

func (m mockApiError) HttpStatusCode() int {
	return m.code
}
