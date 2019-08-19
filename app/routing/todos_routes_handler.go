package routing

import (
	"github.com/lloydmeta/todddo-openapi/internal/api/controllers"
	"github.com/lloydmeta/todddo-openapi/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lloydmeta/todddo-openapi/internal/api/models"
)

type TodosRoutesHandler struct {
	Controller controllers.TodoController
}

// RegisterRoutes takes the given gin.Engine reference and adds the
// routes that it knows how to take care of
func (h *TodosRoutesHandler) RegisterRoutes(ginEngine *gin.Engine) {
	ginEngine.POST("/tasks", h.create)
	ginEngine.GET("/tasks/:id", h.get)
	ginEngine.GET("/tasks", h.list)
	ginEngine.PUT("/tasks/:id", h.update)
	ginEngine.DELETE("/tasks/:id", h.delete)
}

// @Summary Add a new Todo
// @ID create-todo
// @Description Creates a new Todo
// @Accept  json
// @Produce  json
// @Param   todo body models.TodoData true "The request body"
// @Success 200 {object} models.Todo
// @Failure 400 {object} models.Error "Task cannot be empty"
// @Router /tasks [post]
func (h *TodosRoutesHandler) create(c *gin.Context) {
	var apiNewTodo models.TodoData
	if err := c.ShouldBindJSON(&apiNewTodo); err != nil {
		errResp := models.Error{Message: err.Error()}
		c.JSON(http.StatusBadRequest, errResp)
		return
	} else {
		if todo, err := h.Controller.Create(&apiNewTodo); err == nil {
			c.JSON(http.StatusCreated, todo)
		} else {
			c.JSON(err.HttpStatusCode(), err.AsModel())
		}
	}
}

// @Summary Get a Todo by id
// @ID get-existing-todo
// @Description Retrieves a persisted Todo
// @Accept  json
// @Produce  json
// @Param   id path int true "The id of the todo you want to retrieve"
// @Success 200 {object} models.Todo
// @Failure 404 {object} models.Error "Task does not exist"
// @Router /tasks/{id} [get]
func (h *TodosRoutesHandler) get(c *gin.Context) {
	var idPathParam todoIdPathParam
	if err := c.ShouldBindUri(&idPathParam); err != nil {
		errResp := models.Error{Message: err.Error()}
		c.JSON(http.StatusBadRequest, errResp)
		return
	} else {
		id := idPathParam.ID()
		if todo, err := h.Controller.Get(&id); err == nil {
			c.JSON(http.StatusOK, todo)
		} else {
			c.JSON(err.HttpStatusCode(), err.AsModel())
		}
	}
}

// @Summary List all existing Todos
// @ID list-existing-todos
// @Description Retrieves all persisted Todos
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Todo
// @Router /tasks [get]
func (h *TodosRoutesHandler) list(c *gin.Context) {
	list := h.Controller.List()
	c.JSON(http.StatusOK, list)
}

// @Summary Update an existing Todo
// @ID update-todo
// @Description Updates an existing Todo
// @Accept  json
// @Produce  json
// @Param   todo body models.TodoData true "The request body"
// @Param   id path int true "The id of the todo you want to update"
// @Success 200 {object} models.Todo
// @Failure 404 {object} models.Error "Task does not exist"
// @Failure 400 {object} models.Error "Task cannot be empty"
// @Router /tasks/{id} [put]
func (h *TodosRoutesHandler) update(c *gin.Context) {
	var idPathParam todoIdPathParam
	if err := c.ShouldBindUri(&idPathParam); err != nil {
		errResp := models.Error{Message: err.Error()}
		c.JSON(http.StatusBadRequest, errResp)
		return
	} else {
		var apiTodoData models.TodoData
		if err := c.ShouldBindJSON(&apiTodoData); err != nil {
			errResp := models.Error{Message: err.Error()}
			c.JSON(http.StatusBadRequest, errResp)
			return
		} else {
			apiTodo := models.Todo{
				ID:   idPathParam.ID(),
				Task: apiTodoData.Task,
			}
			if todo, err := h.Controller.Update(&apiTodo); err == nil {
				c.JSON(http.StatusOK, todo)
			} else {
				c.JSON(err.HttpStatusCode(), err.AsModel())
			}
		}
	}
}

// @Summary Delete an existing Todo
// @ID delete-todo
// @Description Deletes an existing Todo
// @Accept  json
// @Produce  json
// @Param   id path int true "The id of the todo you want to delete"
// @Success 200 {object} models.Success
// @Failure 404 {object} models.Error "Task does not exist"
// @Router /tasks/{id} [delete]
func (h *TodosRoutesHandler) delete(c *gin.Context) {
	var idPathParam todoIdPathParam
	if err := c.ShouldBindUri(&idPathParam); err != nil {
		errResp := models.Error{Message: err.Error()}
		c.JSON(http.StatusBadRequest, errResp)
		return
	} else {
		id := idPathParam.ID()
		if todo, err := h.Controller.Delete(&id); err == nil {
			c.JSON(http.StatusOK, todo)
		} else {
			c.JSON(err.HttpStatusCode(), err.AsModel())
		}
	}
}

type todoIdPathParam struct {
	UintId uint `uri:"id" binding:"required"`
}

func (t *todoIdPathParam) ID() domain.TodoID {
	return domain.TodoID(t.UintId)
}
