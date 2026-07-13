package scheduling

import (
	"testing"
	"time"

	"github.com/chop1k/medods-test/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

func TestScheduleWeeklyTasks(t *testing.T) {
	templates := []models.Template{
		{
			ID: 1,
			TemplateBody: models.TemplateBody{
				Name:        "Plain weekly task",
				Description: nil,
				Scheduling: &models.Scheduling{
					Type:    models.SchedulingWeekly,
					Include: []string{"wednesday"},
				},
			},
		},
	}

	taskId := 1

	date1 := "07-01-2026"
	date2 := "14-01-2026"
	date3 := "21-01-2026"
	date4 := "28-01-2026"

	result := []models.TaskBody{
		{
			TemplateID: &taskId,
			Status:     models.TaskStatusPending,
			Notes:      nil,
			Date:       &date1,
			StartedAt:  nil,
			EndedAt:    nil,
		},
		{
			TemplateID: &taskId,
			Status:     models.TaskStatusPending,
			Notes:      nil,
			Date:       &date2,
			StartedAt:  nil,
			EndedAt:    nil,
		},
		{
			TemplateID: &taskId,
			Status:     models.TaskStatusPending,
			Notes:      nil,
			Date:       &date3,
			StartedAt:  nil,
			EndedAt:    nil,
		},
		{
			TemplateID: &taskId,
			Status:     models.TaskStatusPending,
			Notes:      nil,
			Date:       &date4,
			StartedAt:  nil,
			EndedAt:    nil,
		},
	}

	from := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC)

	for _, template := range templates {
		t.Run(template.Name, func(t *testing.T) {
			calendar := ScheduleWeeklyTask(template, from, to)

			assert.Equal(t, result, calendar)
		})
	}
}
