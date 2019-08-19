package inmem

import (
	"fmt"
	"testing"

	"github.com/icrowley/fake"

	"github.com/lloydmeta/todddo-openapi/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	repo := MkRepo()
	newTodo := domain.NewTodo{Task: "clean up after yourself"}
	created := repo.Create(&newTodo)
	assert.Equal(t, newTodo.Task, created.Task)
}

func TestGetPresent(t *testing.T) {
	repo := MkRepo()
	newTodo := domain.NewTodo{Task: "clean up after yourself"}
	created := repo.Create(&newTodo)
	retrieved, _ := repo.Get(&created.ID)
	assert.Equal(t, newTodo.Task, retrieved.Task)
}

func TestGetAbsent(t *testing.T) {
	repo := MkRepo()
	id := domain.TodoID(999999)
	_, err := repo.Get(&id)
	assert.Equal(t, true, err != nil)
}

func TestList(t *testing.T) {
	repo := MkRepo()
	toMake := 10
	var createds []domain.Todo
	for i := 0; i < toMake; i++ {
		newTodo := domain.NewTodo{Task: fake.Sentence()}
		createds = append(createds, repo.Create(&newTodo))
	}
	listed := repo.List()
	assert.Equal(t, toMake, len(listed))
	for _, created := range createds {
		var foundInList *domain.Todo
		for _, listed := range listed {
			if listed.ID == created.ID {
				foundInList = &listed
				break
			}
		}
		if foundInList == nil {
			assert.Fail(t, fmt.Sprintf("Not found: [%v]", +created.ID))
		} else {
			assert.Equal(t, created.Task, foundInList.Task, fmt.Sprintf("createds: [%v], listed: [%v]", createds, listed))
		}
	}
}

func TestDeletePresent(t *testing.T) {
	repo := MkRepo()
	newTodo := domain.NewTodo{Task: "clean up after yourself"}
	created := repo.Create(&newTodo)
	deleted, _ := repo.Delete(&created.ID)
	assert.True(t, deleted)

	_, err := repo.Get(&created.ID)
	assert.True(t, err != nil)
}

func TestDeleteAbsent(t *testing.T) {
	repo := MkRepo()
	id := domain.TodoID(99999999)
	deleted, err := repo.Delete(&id)
	assert.False(t, deleted)
	assert.True(t, err != nil)
}

func TestUpdatePresent(t *testing.T) {
	repo := MkRepo()
	newTodo := domain.NewTodo{Task: "clean up after yourself"}
	created := repo.Create(&newTodo)
	created.Task = "do the dishes"
	_, err := repo.Update(&created)
	assert.True(t, err == nil)
	retrieved, _ := repo.Get(&created.ID)
	assert.Equal(t, created.Task, retrieved.Task)
}

func TestUpdateAbsent(t *testing.T) {
	repo := MkRepo()
	update := domain.Todo{
		ID:   domain.TodoID(1235135151),
		Task: "something something",
	}
	_, err := repo.Update(&update)
	assert.True(t, err != nil)
}
