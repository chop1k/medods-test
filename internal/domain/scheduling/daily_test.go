package scheduling

import (
	"testing"
	"time"

	"github.com/chop1k/medods-test/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

func TestScheduleDailyTasks(t *testing.T) {
	templates := []models.Template{
		{
			ID: 1,
			TemplateBody: models.TemplateBody{
				Name:        "Plain daily task",
				Description: nil,
				Scheduling: &models.Scheduling{
					Type: models.SchedulingDaily,
				},
			},
		},
	}

	taskId := 1

	date1 := "01-01-2026"
	date2 := "02-01-2026"
	date3 := "03-01-2026"
	date4 := "04-01-2026"
	date5 := "05-01-2026"
	date6 := "06-01-2026"
	date7 := "07-01-2026"

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
		{
			TemplateID: &taskId,
			Status:     models.TaskStatusPending,
			Notes:      nil,
			Date:       &date5,
			StartedAt:  nil,
			EndedAt:    nil,
		},
		{
			TemplateID: &taskId,
			Status:     models.TaskStatusPending,
			Notes:      nil,
			Date:       &date6,
			StartedAt:  nil,
			EndedAt:    nil,
		},
		{
			TemplateID: &taskId,
			Status:     models.TaskStatusPending,
			Notes:      nil,
			Date:       &date7,
			StartedAt:  nil,
			EndedAt:    nil,
		},
	}

	from := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, time.January, 8, 0, 0, 0, 0, time.UTC)

	for _, template := range templates {
		t.Run(template.Name, func(t *testing.T) {
			calendar := ScheduleDailyTask(template, from, to)

			assert.Equal(t, result, calendar)
		})
	}
}

func TestScheduleDailyTasks_Exclude(t *testing.T) {
	templates := []models.Template{
		{
			ID: 1,
			TemplateBody: models.TemplateBody{
				Name:        "Exclude wildcard",
				Description: nil,
				Scheduling: &models.Scheduling{
					Type:    models.SchedulingDaily,
					Exclude: []string{"03-**-****"},
				},
			},
		},
		{
			ID: 1,
			TemplateBody: models.TemplateBody{
				Name:        "Exclude",
				Description: nil,
				Scheduling: &models.Scheduling{
					Type:    models.SchedulingDaily,
					Exclude: []string{"03-01-2026"},
				},
			},
		},
	}

	taskId := 1

	date1 := "01-01-2026"
	date2 := "02-01-2026"
	date4 := "04-01-2026"
	date5 := "05-01-2026"
	date6 := "06-01-2026"
	date7 := "07-01-2026"

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
			Date:       &date4,
			StartedAt:  nil,
			EndedAt:    nil,
		},
		{
			TemplateID: &taskId,
			Status:     models.TaskStatusPending,
			Notes:      nil,
			Date:       &date5,
			StartedAt:  nil,
			EndedAt:    nil,
		},
		{
			TemplateID: &taskId,
			Status:     models.TaskStatusPending,
			Notes:      nil,
			Date:       &date6,
			StartedAt:  nil,
			EndedAt:    nil,
		},
		{
			TemplateID: &taskId,
			Status:     models.TaskStatusPending,
			Notes:      nil,
			Date:       &date7,
			StartedAt:  nil,
			EndedAt:    nil,
		},
	}

	from := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, time.January, 8, 0, 0, 0, 0, time.UTC)

	for _, template := range templates {
		t.Run(template.Name, func(t *testing.T) {
			calendar := ScheduleDailyTask(template, from, to)

			assert.Equal(t, result, calendar)
		})
	}
}
