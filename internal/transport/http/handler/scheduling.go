package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/chop1k/medods-test/internal/domain/models"
	"github.com/chop1k/medods-test/internal/domain/scheduling"
	"github.com/chop1k/medods-test/internal/repository"
)

type SchedulingHandler struct {
	templatesRepository *repository.TemplatesStorage
	tasksRepository     *repository.TasksStorage
}

func NewSchedulingHandler(storage *repository.TemplatesStorage, tasksStorage *repository.TasksStorage) *SchedulingHandler {
	return &SchedulingHandler{
		templatesRepository: storage,
		tasksRepository:     tasksStorage,
	}
}

func (h *SchedulingHandler) ConnectivityTest(c *gin.Context) {
	c.JSON(http.StatusOK, models.ConnectivityCheck{
		IP:     c.ClientIP(),
		Scheme: c.Request.URL.Scheme,
	})
}

func (h *SchedulingHandler) updateOverdueTasks() {
}

func (h *SchedulingHandler) DailyCronHook(c *gin.Context) {
	h.updateOverdueTasks()

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
			tasks = append(tasks, scheduling.ScheduleDailyTask(template, from, to)...)
		}

		page++
	}

	ids := []int{}

	for _, task := range tasks {
		id, err := h.tasksRepository.Create(nil, task)

		if err != nil {
			panic(err)
		}

		ids = append(ids, id)
	}

	results := []models.Task{}

	for i, id := range ids {
		results = append(results, models.Task{
			ID:       id,
			TaskBody: tasks[i],
		})
	}

	c.JSON(http.StatusOK, results)
}

func (h *SchedulingHandler) WeeklyCronHook(c *gin.Context) {
	page := 1

	tasks := []models.TaskBody{}

	from := time.Now()
	to := from.AddDate(0, 0, 7)

	for {
		templates, err := h.templatesRepository.GetAllWeekly(page, 256)

		if err != nil {
			panic(err)
		}

		if len(templates) == 0 {
			break
		}

		for _, template := range templates {
			tasks = append(tasks, scheduling.ScheduleWeeklyTask(template, from, to)...)
		}

		page++
	}

	ids := []int{}

	for _, task := range tasks {
		id, err := h.tasksRepository.Create(nil, task)

		if err != nil {
			panic(err)
		}

		ids = append(ids, id)
	}

	results := []models.Task{}

	for i, id := range ids {
		results = append(results, models.Task{
			ID:       id,
			TaskBody: tasks[i],
		})
	}

	c.JSON(http.StatusOK, results)
}

func (h *SchedulingHandler) MonthlyCronHook(c *gin.Context) {
	page := 1

	tasks := []models.TaskBody{}

	from := time.Now()
	to := from.AddDate(0, 1, 0)

	for {
		templates, err := h.templatesRepository.GetAllMonthly(page, 256)

		if err != nil {
			panic(err)
		}

		if len(templates) == 0 {
			break
		}

		for _, template := range templates {
			tasks = append(tasks, scheduling.ScheduleMonthlyTask(template, from, to)...)
		}

		page++
	}

	ids := []int{}

	for _, task := range tasks {
		id, err := h.tasksRepository.Create(nil, task)

		if err != nil {
			panic(err)
		}

		ids = append(ids, id)
	}

	results := []models.Task{}

	for i, id := range ids {
		results = append(results, models.Task{
			ID:       id,
			TaskBody: tasks[i],
		})
	}

	c.JSON(http.StatusOK, results)
}

func (h *SchedulingHandler) GetCalendar(c *gin.Context) {
}
