package models

import (
	"github.com/lloydmeta/todddo-openapi/internal/domain"
)

// TodoData models the payload for creating a new Todo
type TodoData struct {
	Task string `json:"task" binding:"required" example:"Buy milk and eggs"`
}

// Todo models the payload for an existing Todo
type Todo struct {
	ID   domain.TodoID `json:"id" binding:"required" example:"1"`
	Task string        `json:"task" binding:"required" example:"Buy milk and eggs"`
}
