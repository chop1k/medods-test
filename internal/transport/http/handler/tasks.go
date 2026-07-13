package handler

import (
	"database/sql"
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/chop1k/medods-test/internal/domain/models"
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

	tasks, count, err := h.repository.GetAll(nil, query.Page, query.Limit, query.Sort, query.SortField)

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, models.TaskListResponse{
		Data: tasks,
		Meta: models.PaginationMeta{
			Total:      len(tasks),
			Page:       query.Page,
			Limit:      query.Limit,
			TotalPages: int(math.Ceil(float64(count) / float64(query.Limit))),
		},
	})
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var body models.TaskBody

	if err := c.ShouldBindJSON(&body); err != nil {
		ValidationError(c, err)

		return
	}

	id, err := h.repository.Create(nil, body)

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

	task, err := h.repository.GetById(nil, param.TaskID)

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

	task, err := h.repository.GetById(nil, param.TaskID)

	if err != nil {
		panic(err)
	}

	if task.DeletedAt != nil {
		NotFound(c, "not found")

		return
	}

	var runningTask models.RunningTaskBody

	if err := c.ShouldBindBodyWithJSON(&runningTask); err == nil {
		oldTask, err := h.repository.GetById(nil, param.TaskID)

		if err != nil {
			panic(err)
		}

		if oldTask.Status != models.TaskStatusPending {
			ValidationError(c, errors.New("task should be pending"))

			return
		}

		task := models.TaskBody{
			TemplateID: oldTask.TemplateID,
			Status:     models.TaskStatusRunning,
			Notes:      oldTask.Notes,
			Date:       oldTask.Date,
			StartedAt:  runningTask.StartedAt,
			EndedAt:    oldTask.EndedAt,
		}

		updated, err := h.repository.UpdateById(nil, param.TaskID, task)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, updated)

		return
	}

	var finishedTask models.FinishedTaskBody

	if err := c.ShouldBindBodyWithJSON(&finishedTask); err == nil {
		oldTask, err := h.repository.GetById(nil, param.TaskID)

		if err != nil {
			panic(err)
		}

		if oldTask.Status != models.TaskStatusRunning {
			ValidationError(c, errors.New("task should be running"))

			return
		}

		task := models.TaskBody{
			TemplateID: oldTask.TemplateID,
			Status:     models.TaskStatusFinished,
			Notes:      oldTask.Notes,
			Date:       oldTask.Date,
			StartedAt:  oldTask.StartedAt,
			EndedAt:    finishedTask.EndedAt,
		}

		updated, err := h.repository.UpdateById(nil, param.TaskID, task)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, updated)

		return
	}

	var cancelledTask models.CancelledTaskBody

	if err := c.ShouldBindBodyWithJSON(&cancelledTask); err == nil {
		oldTask, err := h.repository.GetById(nil, param.TaskID)

		if err != nil {
			panic(err)
		}

		var task models.TaskBody

		switch oldTask.Status {
		case models.TaskStatusPending:
			task = models.TaskBody{
				TemplateID: oldTask.TemplateID,
				Status:     models.TaskStatusCancelled,
				Notes:      oldTask.Notes,
				Date:       oldTask.Date,
				StartedAt:  nil,
				EndedAt:    nil,
			}
		case models.TaskStatusRunning:
			task = models.TaskBody{
				TemplateID: oldTask.TemplateID,
				Status:     models.TaskStatusCancelled,
				Notes:      oldTask.Notes,
				Date:       oldTask.Date,
				StartedAt:  oldTask.StartedAt,
				EndedAt:    cancelledTask.EndedAt,
			}
		default:
			ValidationError(c, errors.New("task should be either running or pending"))

			return
		}

		updated, err := h.repository.UpdateById(nil, param.TaskID, task)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, updated)

		return
	}

	var movedTask models.MovedTaskBody

	if err := c.ShouldBindBodyWithJSON(&movedTask); err == nil {
		oldTask, err := h.repository.GetById(nil, param.TaskID)

		if err != nil {
			panic(err)
		}

		if oldTask.Status != models.TaskStatusPending {
			ValidationError(c, errors.New("task should be pending"))

			return
		}

		updated, err := h.repository.Transaction(func(tx *sql.Tx) (any, error) {
			newTask := models.TaskBody{
				TemplateID: nil,
				MovedId:    nil,
				Status:     models.TaskStatusPending,
				Date:       movedTask.Date,
				StartedAt:  nil,
				EndedAt:    nil,
			}

			created, err := h.repository.Create(tx, newTask)

			if err != nil {
				return nil, err
			}

			task := models.TaskBody{
				TemplateID: oldTask.TemplateID,
				MovedId:    &created,
				Status:     models.TaskStatusMoved,
				Notes:      oldTask.Notes,
				Date:       oldTask.Date,
				StartedAt:  oldTask.StartedAt,
				EndedAt:    oldTask.EndedAt,
			}

			return h.repository.UpdateById(tx, param.TaskID, task)
		})

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, updated)

		return
	}

	ValidationError(c, errors.New("unknown body format"))
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	var param models.TaskIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		ValidationError(c, err)

		return
	}

	task, err := h.repository.GetById(nil, param.TaskID)

	if err != nil {
		panic(err)
	}

	if task.DeletedAt != nil {
		NotFound(c, "not found")

		return
	}

	err = h.repository.RemoveById(nil, param.TaskID)

	if err != nil {
		panic(err)
	}

	c.Status(http.StatusNoContent)
}
