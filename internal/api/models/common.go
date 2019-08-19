package models

// Error models errors as JSON in the API
type Error struct {
	Message string `json:"message" binding:"required" example:"Something went wrong :("`
}

type Success struct {
	Message string `json:"message" binding:"required" example:"Something went right!"`
}

type ApiError interface {
	error
	AsModel() Error
	HttpStatusCode() int
}
