package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/chop1k/medods-test/internal/database"
	"github.com/chop1k/medods-test/internal/models"
	"github.com/chop1k/medods-test/internal/scheduling"
)

type SchedulingHandler struct {
	templatesRepository *database.TemplatesStorage
	tasksRepository     *database.TasksStorage
}

func NewSchedulingHandler(storage *database.TemplatesStorage) *SchedulingHandler {
	return &SchedulingHandler{
		templatesRepository: storage,
	}
}

func (h *SchedulingHandler) ConnectivityTest(c *gin.Context) {
	c.JSON(http.StatusOK, models.ConnectivityCheck{
		IP:     c.ClientIP(),
		Scheme: c.Request.URL.Scheme,
	})
}

func (h *SchedulingHandler) DailyCronHook(c *gin.Context) {
	page := 1

	tasks := []models.TaskBody{}

	from := time.Now()
	to := from.AddDate(0, 0, 1)

	for {
		templates, err := h.templatesRepository.GetAllDaily(page, 256)

		if err != nil {
			panic(err)
		}

		if len(templates) == 0 {
			break
		}

		for _, template := range templates {
			_tasks, err := scheduling.ScheduleTemplate(template, from, to)

			if err != nil {
				panic(err)
			}

			tasks = append(tasks, _tasks...)
		}
	}

	ids, err := h.tasksRepository.CreateBulk(tasks)

	if err != nil {
		panic(err)
	}

	results := []models.Task{}

	for i, id := range ids {
		results[i] = models.Task{
			ID:       id,
			TaskBody: tasks[i],
		}
	}

	c.JSON(http.StatusOK, results)
}

func (h *SchedulingHandler) WeeklyCronHook(c *gin.Context) {
}

func (h *SchedulingHandler) MonthlyCronHook(c *gin.Context) {
}

func (h *SchedulingHandler) GetCalendar(c *gin.Context) {
}
