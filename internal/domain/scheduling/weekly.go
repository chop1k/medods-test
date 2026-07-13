package scheduling

import (
	"time"

	"github.com/chop1k/medods-test/internal/domain/models"
)

func ScheduleWeeklyTask(template models.Template, from time.Time, to time.Time) []models.TaskBody {
	from = trimTime(from)
	to = trimTime(to)

	tasks := []models.TaskBody{}

	for i := 0; i < diffDays(from, to); i++ {
		for _, day := range template.Scheduling.Include {
			date := from.AddDate(0, 0, i)

			if weekday(day) != date.Weekday() {
				continue
			}

			format := date.Format(dateFormat)

			tasks = append(tasks, models.TaskBody{
				TemplateID: &template.ID,
				Status:     models.TaskStatusPending,
				Notes:      nil,
				Date:       &format,
				StartedAt:  nil,
				EndedAt:    nil,
			})
		}
	}

	return tasks
}
