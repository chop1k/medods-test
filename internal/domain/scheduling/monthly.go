package scheduling

import (
	"strconv"
	"time"

	"github.com/chop1k/medods-test/internal/domain/models"
)

func dayOfMonth(day string) int {
	value, err := strconv.Atoi(day)

	if err != nil {
		return -1
	}

	return value
}

func ScheduleMonthlyTask(template models.Template, from time.Time, to time.Time) []models.TaskBody {
	from = trimTime(from)
	to = trimTime(to)

	tasks := []models.TaskBody{}

	for i := 1; i <= diffDays(from, to); i++ {
		for _, day := range template.Scheduling.Include {
			date := from.AddDate(0, 0, i)

			if weekday(day) != date.Weekday() && dayOfMonth(day) != date.Day() {
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
