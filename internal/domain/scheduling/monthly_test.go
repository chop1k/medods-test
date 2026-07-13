package scheduling

import (
	"testing"
	"time"

	"github.com/chop1k/medods-test/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

func TestScheduleMonthlyTasks(t *testing.T) {
	templates := []models.Template{
		{
			ID: 1,
			TemplateBody: models.TemplateBody{
				Name:        "Plain monthly task",
				Description: nil,
				Scheduling: &models.Scheduling{
					Type:    models.SchedulingMonthly,
					Include: []string{"22"},
				},
			},
		},
	}

	taskId := 1

	date := "22-01-2026"

	result := []models.TaskBody{
		{
			TemplateID: &taskId,
			Status:     models.TaskStatusPending,
			Notes:      nil,
			Date:       &date,
			StartedAt:  nil,
			EndedAt:    nil,
		},
	}

	from := time.Date(2026, time.January, 0, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, time.February, 0, 0, 0, 0, 0, time.UTC)

	for _, template := range templates {
		t.Run(template.Name, func(t *testing.T) {
			calendar := ScheduleMonthlyTask(template, from, to)

			assert.Equal(t, result, calendar)
		})
	}
}
