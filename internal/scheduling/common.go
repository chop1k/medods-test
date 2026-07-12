package scheduling

import (
	"fmt"
	"time"

	"github.com/chop1k/medods-test/internal/models"
	"github.com/emirpasic/gods/lists/arraylist"
)

const DateFormat = "02-01-2006"

func ScheduleTemplate(template models.Template, from time.Time, to time.Time) ([]models.TaskBody, error) {
	switch template.Scheduling.Type {
	case models.SchedulingDaily:
		return scheduleDailyTask(template, from, to)
	case models.SchedulingWeekly:
		return scheduleWeeklyTask(template, from, to)
	case models.SchedulingMonthly:
		return scheduleMonthlyTask(template, from, to)
	default:
		return nil, fmt.Errorf("unknown scheduling type %s", template.Scheduling.Type)
	}
}

func scheduleDailyTask(templates models.Template, from time.Time, to time.Time) ([]models.TaskBody, error) {
	// check if day excluded
	// make a task

	from = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location())
	to = time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, to.Location())

	diff := to.Sub(from)
	days := int(diff.Hours() / 24)

	exclude := arraylist.New(templates.Scheduling.Exclude)

	tasks := []models.TaskBody{}

	for i := 1; i <= days; i++ {
		date := from.AddDate(0, 0, days)

		if exclude.Empty() {
			continue
		}

		if exclude.Contains(date.Format(DateFormat)) {
			continue
		}

		tasks = append(tasks, models.TaskBody{
			TemplateID: templates.ID,
			Status:     models.TaskStatusPending,
			Notes:      nil,
			StartedAt:  nil,
			EndedAt:    nil,
		})
	}

	return tasks, nil
}

func scheduleWeeklyTask(templates models.Template, from time.Time, to time.Time) ([]models.TaskBody, error) {
	panic("not implemented")
}

func scheduleMonthlyTask(templates models.Template, from time.Time, to time.Time) ([]models.TaskBody, error) {
	panic("not implemented")
}

func scheduleOneshotTask(templates models.Template, from time.Time, to time.Time) ([]models.TaskBody, error) {
	panic("not implemented")
}
