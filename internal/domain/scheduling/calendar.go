package scheduling

import "github.com/chop1k/medods-test/internal/domain/models"

type Calendar map[string][]models.TaskBody

func NewCalendarFromTasks(tasks []models.TaskBody) Calendar {
	calendar := make(Calendar)

	for _, task := range tasks {
		date := *task.Date

		days, ok := calendar[date]

		if !ok {
			days = []models.TaskBody{}
		}

		days = append(days, task)

		calendar[date] = days
	}

	return calendar
}

func (c Calendar) Merge(calendar Calendar) {
	for k, v := range calendar {
		tasks := c[k]

		c[k] = append(tasks, v...)
	}
}
