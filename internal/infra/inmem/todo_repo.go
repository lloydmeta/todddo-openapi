package inmem

import (
	"sort"
	"sync"

	"github.com/lloydmeta/todddo-openapi/internal/domain"
)

type repoImpl struct {
	mutex  sync.Mutex
	lastId domain.TodoID
	stored map[domain.TodoID]persistedTask
}

type persistedTask struct {
	task string
}

// MkRepo returns a new TodoRepo based on an in-mem implementation
func MkRepo() domain.TodoRepo {
	return &repoImpl{
		stored: make(map[domain.TodoID]persistedTask),
	}
}

func (r *repoImpl) Create(newTodo *domain.NewTodo) domain.Todo {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	id := r.lastId + 1
	r.lastId = id
	persisted := persistedTask{newTodo.Task}
	r.stored[id] = persisted
	return domain.Todo{
		ID:   id,
		Task: newTodo.Task,
	}
}

func (r *repoImpl) Get(id *domain.TodoID) (domain.Todo, domain.TodoRepoError) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if retrieved, exists := r.stored[*id]; exists {
		return domain.Todo{
			ID:   *id,
			Task: retrieved.task,
		}, nil
	} else {
		return domain.Todo{}, domain.TodoNotFound{ID: *id}
	}

}
func (r *repoImpl) List() []domain.Todo {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	retrieved := make([]domain.Todo, 0, len(r.stored))
	for id, v := range r.stored {
		retrieved = append(retrieved, domain.Todo{
			ID:   id,
			Task: v.task,
		})
	}
	sort.SliceStable(retrieved, func(i, j int) bool { return retrieved[i].ID < retrieved[j].ID })
	return retrieved
}

func (r *repoImpl) Delete(id *domain.TodoID) (bool, domain.TodoRepoError) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if _, exists := r.stored[*id]; exists {
		delete(r.stored, *id)
		return true, nil
	} else {
		return false, domain.TodoNotFound{ID: *id}
	}
}

func (r *repoImpl) Update(todo *domain.Todo) (domain.Todo, domain.TodoRepoError) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if _, exists := r.stored[todo.ID]; exists {
		r.stored[todo.ID] = persistedTask{task: todo.Task}
		return *todo, nil
	} else {
		return domain.Todo{}, domain.TodoNotFound{ID: todo.ID}
	}
}
