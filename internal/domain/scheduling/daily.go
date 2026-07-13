package scheduling

import (
	"time"

	"github.com/chop1k/medods-test/internal/domain/models"
)

func ScheduleDailyTask(templates models.Template, from time.Time, to time.Time) []models.TaskBody {
	from = trimTime(from)
	to = trimTime(to)

	tasks := []models.TaskBody{}

	for i := 0; i < diffDays(from, to); i++ {
		date := from.AddDate(0, 0, i).Format(dateFormat)

		tasks = append(tasks, models.TaskBody{
			TemplateID: &templates.ID,
			Status:     models.TaskStatusPending,
			Notes:      nil,
			Date:       &date,
			StartedAt:  nil,
			EndedAt:    nil,
		})
	}

	filteredTasks := []models.TaskBody{}

	for _, task := range tasks {
		date := *task.Date
		var filter bool

		for _, exclude := range templates.Scheduling.Exclude {
			if dateMatched(date, exclude) {
				filter = true

				break
			}
		}

		if filter {
			continue
		}

		filteredTasks = append(filteredTasks, task)
	}

	return filteredTasks
}
