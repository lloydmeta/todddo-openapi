package app

import (
	"github.com/lloydmeta/todddo-openapi/internal/api/controllers"
	"github.com/lloydmeta/todddo-openapi/internal/domain"
	"github.com/lloydmeta/todddo-openapi/internal/domain/services"
	"github.com/lloydmeta/todddo-openapi/internal/infra/inmem"
)

type Components struct {
	Controllers Controllers
	Services    Services
	Repos       Repos
}

// MkDefaultComponents returns default components
func MkDefaultComponents() Components {
	repoComponents := Repos{TodoRepo: inmem.MkRepo()}
	serviceComponents := Services{TodoService: services.MkTodoService(repoComponents.TodoRepo)}
	controllerComponents := Controllers{TodoController: controllers.MkTodosController(serviceComponents.TodoService)}
	return Components{
		Controllers: controllerComponents,
		Services:    serviceComponents,
		Repos:       repoComponents,
	}
}

type Controllers struct {
	TodoController controllers.TodoController
}

type Services struct {
	TodoService services.TodoService
}

type Repos struct {
	TodoRepo domain.TodoRepo
}
