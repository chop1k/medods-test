package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/chop1k/medods-test/internal/models"
	"github.com/chop1k/medods-test/internal/responses"
)

type TaskHandler struct {
}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{}
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	var query models.ListTasksQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		responses.ValidationError(c, err)
		return
	}

	// TODO: fetch paginated tasks from the database using query.Page,
	// query.Limit, query.Sort and query.SortField.

	c.JSON(http.StatusOK, models.TaskListResponse{
		Data: []models.Task{},
		Meta: models.PaginationMeta{
			Total:      0,
			Page:       query.Page,
			Limit:      query.Limit,
			TotalPages: 0,
		},
	})
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var body models.TaskBody
	if err := c.ShouldBindJSON(&body); err != nil {
		responses.ValidationError(c, err)
		return
	}

	// TODO: persist the new task and obtain its generated ID.
	created := models.Task{
		ID:       0,
		TaskBody: body,
	}

	// TODO: set the Location header to the URL of the created resource,
	// e.g. c.Header("Location", fmt.Sprintf("/tasks/%d", created.ID))

	c.JSON(http.StatusCreated, created)
}

func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	var param models.TaskIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		responses.ValidationError(c, err)
		return
	}

	// TODO: fetch the task by param.TaskID from the database.
	// If it does not exist, respond with httpresponse.NotFound(c, "...").

	c.JSON(http.StatusOK, models.Task{ID: param.TaskID})
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	var param models.TaskIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		responses.ValidationError(c, err)
		return
	}

	var body models.TaskBody
	if err := c.ShouldBindJSON(&body); err != nil {
		responses.ValidationError(c, err)
		return
	}

	// TODO: verify the task with param.TaskID exists
	// (httpresponse.NotFound(c, "...") if not), then persist the update.
	updated := models.Task{
		ID:       param.TaskID,
		TaskBody: body,
	}

	c.JSON(http.StatusOK, updated)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	var param models.TaskIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		responses.ValidationError(c, err)
		return
	}

	// TODO: verify the task with param.TaskID exists
	// (httpresponse.NotFound(c, "...") if not), then delete it.

	c.Status(http.StatusNoContent)
}
