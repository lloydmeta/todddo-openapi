package services

import (
	"testing"

	"github.com/lloydmeta/todddo-openapi/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreateValidData(t *testing.T) {
	mockRepo := mockRepo{}
	mockRepo.create = func(newTodo *domain.NewTodo) domain.Todo {
		return domain.Todo{
			ID:   domain.TodoID(123),
			Task: newTodo.Task,
		}
	}
	service := todoServiceImpl{Repo: &mockRepo}
	newTodo := domain.NewTodo{Task: "do something"}
	_, err := service.Create(&newTodo)
	assert.Equal(t, uint(1), mockRepo.createCalled)
	assert.True(t, err == nil)
}

func TestCreateInvalidData(t *testing.T) {
	mockRepo := mockRepo{}
	service := todoServiceImpl{Repo: &mockRepo}
	newTodo := domain.NewTodo{Task: ""}
	_, err := service.Create(&newTodo)
	assert.Equal(t, uint(0), mockRepo.createCalled)
	assert.True(t, err != nil)
}

func TestUpdateValidData(t *testing.T) {
	mockRepo := mockRepo{}
	mockRepo.update = func(todo *domain.Todo) (domain.Todo, domain.TodoRepoError) {
		return *todo, nil
	}
	service := todoServiceImpl{Repo: &mockRepo}
	updatedTodo := domain.Todo{ID: domain.TodoID(123), Task: "do something"}
	_, err := service.Update(&updatedTodo)
	assert.Equal(t, uint(1), mockRepo.updateCalled)
	assert.True(t, err == nil)
}

func TestUpdateInvalidData(t *testing.T) {
	mockRepo := mockRepo{}
	service := todoServiceImpl{Repo: &mockRepo}
	updatedTodo := domain.Todo{ID: domain.TodoID(123), Task: ""}
	_, err := service.Update(&updatedTodo)
	assert.Equal(t, uint(0), mockRepo.updateCalled)
	assert.True(t, err != nil)
}

func TestUpdateNotFound(t *testing.T) {
	mockRepo := mockRepo{}
	stored := false
	mockRepo.update = func(todo *domain.Todo) (domain.Todo, domain.TodoRepoError) {
		stored = true
		return domain.Todo{}, domain.TodoNotFound{ID: todo.ID}
	}
	service := todoServiceImpl{Repo: &mockRepo}
	updatedTodo := domain.Todo{ID: domain.TodoID(123), Task: "hello"}
	_, err := service.Update(&updatedTodo)
	assert.Equal(t, uint(1), mockRepo.updateCalled)
	assert.True(t, err != nil)
	assert.True(t, stored)
}

func TestList(t *testing.T) {
	mockRepo := mockRepo{}
	existing := domain.Todo{ID: domain.TodoID(123), Task: "hello"}
	mockRepo.list = func() []domain.Todo {
		return []domain.Todo{existing}
	}
	service := todoServiceImpl{Repo: &mockRepo}
	listed := service.List()
	assert.Equal(t, uint(1), mockRepo.listCalled)
	assert.ElementsMatch(t, []domain.Todo{existing}, listed)
}

func TestGetOk(t *testing.T) {
	mockRepo := mockRepo{}
	existing := domain.Todo{ID: domain.TodoID(123), Task: "hello"}
	mockRepo.get = func(id *domain.TodoID) (domain.Todo, domain.TodoRepoError) {
		return existing, nil
	}
	service := todoServiceImpl{Repo: &mockRepo}
	id := domain.TodoID(123)
	found, err := service.Get(&id)
	assert.Equal(t, uint(1), mockRepo.getCalled)
	assert.Equal(t, existing, found)
	assert.True(t, err == nil)
}

func TestGetNotFound(t *testing.T) {
	mockRepo := mockRepo{}
	mockRepo.get = func(id *domain.TodoID) (domain.Todo, domain.TodoRepoError) {
		return domain.Todo{}, domain.TodoNotFound{ID: *id}
	}
	service := todoServiceImpl{Repo: &mockRepo}
	id := domain.TodoID(123)
	_, err := service.Get(&id)
	assert.Equal(t, uint(1), mockRepo.getCalled)
	assert.True(t, err != nil)
}

func TestDeleteOk(t *testing.T) {
	mockRepo := mockRepo{}
	mockRepo.delete = func(id *domain.TodoID) (bool, domain.TodoRepoError) {
		return true, nil
	}
	service := todoServiceImpl{Repo: &mockRepo}
	id := domain.TodoID(123)
	deleted, err := service.Delete(&id)
	assert.True(t, deleted)
	assert.True(t, err == nil)
}

func TestDeleteNotFound(t *testing.T) {
	mockRepo := mockRepo{}
	mockRepo.delete = func(id *domain.TodoID) (bool, domain.TodoRepoError) {
		return false, domain.TodoNotFound{ID: *id}
	}
	service := todoServiceImpl{Repo: &mockRepo}
	id := domain.TodoID(123)
	deleted, err := service.Delete(&id)
	assert.False(t, deleted)
	assert.True(t, err != nil)
}

// mocks

type mockRepo struct {
	create       func(newTodo *domain.NewTodo) domain.Todo
	createCalled uint
	get          func(id *domain.TodoID) (domain.Todo, domain.TodoRepoError)
	getCalled    uint
	list         func() []domain.Todo
	listCalled   uint
	delete       func(id *domain.TodoID) (bool, domain.TodoRepoError)
	deleteCalled uint
	update       func(todo *domain.Todo) (domain.Todo, domain.TodoRepoError)
	updateCalled uint
}

func (r *mockRepo) Create(newTodo *domain.NewTodo) domain.Todo {
	defer func() { r.createCalled++ }()
	return r.create(newTodo)
}

func (r *mockRepo) Get(id *domain.TodoID) (domain.Todo, domain.TodoRepoError) {
	defer func() { r.getCalled++ }()
	return r.get(id)
}
func (r *mockRepo) List() []domain.Todo {
	defer func() { r.listCalled++ }()
	return r.list()
}

func (r *mockRepo) Delete(id *domain.TodoID) (bool, domain.TodoRepoError) {
	defer func() { r.deleteCalled++ }()
	return r.delete(id)
}

func (r *mockRepo) Update(todo *domain.Todo) (domain.Todo, domain.TodoRepoError) {
	defer func() { r.updateCalled++ }()
	return r.update(todo)
}
