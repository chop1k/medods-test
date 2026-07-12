package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/chop1k/medods-test/internal/models"
	"github.com/chop1k/medods-test/internal/repository"
)

type TaskHandler struct {
	repository *repository.TasksStorage
}

func NewTaskHandler(repository *repository.TasksStorage) *TaskHandler {
	return &TaskHandler{
		repository: repository,
	}
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	var query models.ListTasksQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ValidationError(c, err)

		return
	}

	tasks, err := h.repository.GetAll(query.Page, query.Limit)

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, models.TaskListResponse{
		Data: tasks,
		Meta: models.PaginationMeta{
			Total:      len(tasks),
			Page:       query.Page,
			Limit:      query.Limit,
			TotalPages: 0,
		},
	})
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var body models.TaskBody

	if err := c.ShouldBindJSON(&body); err != nil {
		ValidationError(c, err)

		return
	}

	id, err := h.repository.Create(body)

	if err != nil {
		panic(err)
	}

	created := models.Task{
		ID:       id,
		TaskBody: body,
	}

	c.Header("Location", "/v1/tasks/task/"+strconv.Itoa(id))
	c.JSON(http.StatusCreated, created)
}

func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	var param models.TaskIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		ValidationError(c, err)

		return
	}

	task, err := h.repository.GetById(param.TaskID)

	if err != nil {
		NotFound(c, err.Error())

		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	var param models.TaskIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		ValidationError(c, err)

		return
	}

	var body models.TaskBody
	if err := c.ShouldBindJSON(&body); err != nil {
		ValidationError(c, err)

		return
	}

	updated, err := h.repository.UpdateById(param.TaskID, body)

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, updated)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	var param models.TaskIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		ValidationError(c, err)

		return
	}

	err := h.repository.RemoveById(param.TaskID)

	if err != nil {
		panic(err)
	}

	c.Status(http.StatusNoContent)
}
