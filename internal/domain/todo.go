package domain

import (
	"fmt"
)

// TodoID is the identifier for a Todo
type TodoID uint64

// NewTodo is for persisting a new Todo
type NewTodo struct {
	Task string
}

// Todo is a persisted Todo
type Todo struct {
	ID   TodoID
	Task string
}

// TodoRepo is an interface for managing the persistence lifecycle
// of a Todo
type TodoRepo interface {
	Create(newTodo *NewTodo) Todo
	Get(id *TodoID) (Todo, TodoRepoError)
	List() []Todo
	Delete(id *TodoID) (bool, TodoRepoError)
	Update(todo *Todo) (Todo, TodoRepoError)
}

// <-- Errors

// TodoRepoError is an error interface for TodoRepo
type TodoRepoError interface {
	error
	Id() TodoID
}

// TodoNotFound is returned when the repo cannot find
// a repo by a given TodoId
type TodoNotFound struct {
	ID TodoID
}

func (e TodoNotFound) Error() string {
	return fmt.Sprintf("Could not find [%v] in in-mem repo", e.ID)
}

func (e TodoNotFound) Id() TodoID {
	return e.ID
}

//     Errors -->
