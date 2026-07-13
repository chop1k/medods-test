package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/chop1k/medods-test/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const dateFormat = "02-01-2006"

func dailyCronURL() string {
	return testURL + "/v1/scheduling/daily-cron-hook"
}

func weeklyCronURL() string {
	return testURL + "/v1/scheduling/weekly-cron-hook"
}

func monthlyCronURL() string {
	return testURL + "/v1/scheduling/monthly-cron-hook"
}

func validDailyTemplates() []models.TemplateBody {
	return []models.TemplateBody{
		{
			Name:        "Every day task",
			Description: nil,
			Enabled:     true,
			StartsAt:    "00:00:00",
			EndsAt:      "01:00:00",
			Scheduling: &models.Scheduling{
				Type: models.SchedulingDaily,
			},
		},
	}
}

func expectedDailyTasks(ids []*int) []models.TaskBody {
	date := time.Now().Format(dateFormat)

	return []models.TaskBody{
		{
			TemplateID: ids[0],
			Status:     models.TaskStatusPending,
			Notes:      nil,
			Date:       &date,
		},
	}
}

func validWeeklyTemplates() []models.TemplateBody {
	return []models.TemplateBody{
		{
			Name:        "Every monday task",
			Description: nil,
			Enabled:     true,
			StartsAt:    "00:00:00",
			EndsAt:      "01:00:00",
			Scheduling: &models.Scheduling{
				Type:    models.SchedulingWeekly,
				Include: []string{"monday"},
			},
		},
	}
}

func expectedWeeklyTasks(ids []*int) []models.TaskBody {
	date := time.Now()

	switch date.Weekday() {
	case time.Tuesday:
		date = date.AddDate(0, 0, 6)
	case time.Wednesday:
		date = date.AddDate(0, 0, 5)
	case time.Thursday:
		date = date.AddDate(0, 0, 4)
	case time.Friday:
		date = date.AddDate(0, 0, 3)
	case time.Saturday:
		date = date.AddDate(0, 0, 2)
	case time.Sunday:
		date = date.AddDate(0, 0, 1)
	}

	formattedDate := date.Format(dateFormat)

	return []models.TaskBody{
		{
			TemplateID: ids[0],
			Status:     models.TaskStatusPending,
			Notes:      nil,
			Date:       &formattedDate,
		},
	}
}

func validMonthlyTemplates() []models.TemplateBody {
	return []models.TemplateBody{
		{
			Name:        "Every day task",
			Description: nil,
			Enabled:     true,
			StartsAt:    "00:00:00",
			EndsAt:      "01:00:00",
			Scheduling: &models.Scheduling{
				Type:    models.SchedulingMonthly,
				Include: []string{"07"},
			},
		},
	}
}

func expectedMonthlyTasks(ids []*int) []models.TaskBody {
	now := time.Now()

	var date string

	if now.Day() > 7 {
		date = time.Date(now.Year(), now.Month()+1, 7, 0, 0, 0, 0, now.Location()).Format(dateFormat)
	} else {
		date = time.Date(now.Year(), now.Month(), 7, 0, 0, 0, 0, now.Location()).Format(dateFormat)
	}

	return []models.TaskBody{
		{
			TemplateID: ids[0],
			Status:     models.TaskStatusPending,
			Notes:      nil,
			Date:       &date,
		},
	}
}

func TestScheduleDailyCronHook(t *testing.T) {
	TruncateDB(t)

	bodies := validDailyTemplates()

	for _, body := range bodies {
		t.Run(body.Name, func(t *testing.T) {
			templateJson, err := json.Marshal(body)
			require.Nil(t, err, "cannot marshal the body", err)

			request, err := http.NewRequest(http.MethodPost, templatesCollectionURL(), bytes.NewReader(templateJson))
			require.Nil(t, err, "cannot create create request", err)

			response, err := testClient.Do(request)
			require.Nil(t, err, "create request failed", err)

			defer response.Body.Close()

			require.Equal(t, http.StatusCreated, response.StatusCode)
		})
	}

	getCollectionRequest, err := http.NewRequest(http.MethodGet, templatesCollectionURL(), nil)
	require.Nil(t, err, "cannot create get template request", err)

	getCollectionResponse, err := testClient.Do(getCollectionRequest)
	require.Nil(t, err, "get request failed")
	require.Equal(t, http.StatusOK, getCollectionResponse.StatusCode)
	defer getCollectionResponse.Body.Close()

	var templates models.TemplateListResponse
	err = json.NewDecoder(getCollectionResponse.Body).Decode(&templates)
	require.Nil(t, err, "get templates endpoint returned unknown format", err)

	assert.Equal(t, len(bodies), len(templates.Data))

	for i, template := range templates.Data {
		assert.Equal(t, template.Name, bodies[i].Name)
		assert.Equal(t, template.Description, bodies[i].Description)
		assert.Equal(t, template.StartsAt, bodies[i].StartsAt)
		assert.Equal(t, template.EndsAt, bodies[i].EndsAt)
		assert.Equal(t, template.Scheduling, bodies[i].Scheduling)
	}

	dailyCronRequest, err := http.NewRequest(http.MethodPost, dailyCronURL(), nil)
	require.Nil(t, err, "cannot create daily cron request", err)

	dailyCronResponse, err := testClient.Do(dailyCronRequest)
	require.Nil(t, err, "daily cron request failed")
	require.Equal(t, http.StatusOK, dailyCronResponse.StatusCode)
	defer dailyCronResponse.Body.Close()

	var tasks []models.Task
	err = json.NewDecoder(dailyCronResponse.Body).Decode(&tasks)
	require.Nil(t, err, "daily cron endpoint returned unknown format", err)

	require.Equal(t, 1, len(tasks))

	ids := []*int{}

	for _, task := range tasks {
		ids = append(ids, &task.ID)
	}

	expectedTasks := expectedDailyTasks(ids)

	assert.Equal(t, len(tasks), len(expectedTasks))

	for i, task := range tasks {
		assert.Equal(t, expectedTasks[i].TemplateID, task.TemplateID)
		assert.Equal(t, expectedTasks[i].Status, task.Status)
		assert.Equal(t, expectedTasks[i].Notes, task.Notes)
		assert.Equal(t, expectedTasks[i].StartedAt, task.StartedAt)
		assert.Equal(t, expectedTasks[i].EndedAt, task.EndedAt)
		assert.Equal(t, expectedTasks[i].Date, task.Date)
	}

	getTasksRequest, err := http.NewRequest(http.MethodGet, tasksCollectionURL(), nil)
	require.Nil(t, err, "cannot create get tasks request", err)

	getTasksResponse, err := testClient.Do(getTasksRequest)
	require.Nil(t, err, "get request failed")
	require.Equal(t, http.StatusOK, getTasksResponse.StatusCode)
	defer getTasksResponse.Body.Close()

	var getTasksResult models.TaskListResponse
	err = json.NewDecoder(getTasksResponse.Body).Decode(&getTasksResult)
	require.Nil(t, err, "get tasks endpoint returned unknown format", err)

	assert.Equal(t, tasks, getTasksResult.Data)
}

func TestScheduleWeeklyCronHook(t *testing.T) {
	TruncateDB(t)

	bodies := validWeeklyTemplates()

	for _, body := range bodies {
		t.Run(body.Name, func(t *testing.T) {
			templateJson, err := json.Marshal(body)
			require.Nil(t, err, "cannot marshal the body", err)

			request, err := http.NewRequest(http.MethodPost, templatesCollectionURL(), bytes.NewReader(templateJson))
			require.Nil(t, err, "cannot create create request", err)

			response, err := testClient.Do(request)
			require.Nil(t, err, "create request failed", err)

			defer response.Body.Close()

			require.Equal(t, http.StatusCreated, response.StatusCode)
		})
	}

	getCollectionRequest, err := http.NewRequest(http.MethodGet, templatesCollectionURL(), nil)
	require.Nil(t, err, "cannot create get template request", err)

	getCollectionResponse, err := testClient.Do(getCollectionRequest)
	require.Nil(t, err, "get request failed")
	require.Equal(t, http.StatusOK, getCollectionResponse.StatusCode)
	defer getCollectionResponse.Body.Close()

	var templates models.TemplateListResponse
	err = json.NewDecoder(getCollectionResponse.Body).Decode(&templates)
	require.Nil(t, err, "get templates endpoint returned unknown format", err)

	assert.Equal(t, len(bodies), len(templates.Data))

	for i, template := range templates.Data {
		assert.Equal(t, template.Name, bodies[i].Name)
		assert.Equal(t, template.Description, bodies[i].Description)
		assert.Equal(t, template.StartsAt, bodies[i].StartsAt)
		assert.Equal(t, template.EndsAt, bodies[i].EndsAt)
		assert.Equal(t, template.Scheduling, bodies[i].Scheduling)
	}

	weeklyCronRequest, err := http.NewRequest(http.MethodPost, weeklyCronURL(), nil)
	require.Nil(t, err, "cannot create weekly cron request", err)

	weeklyCronResponse, err := testClient.Do(weeklyCronRequest)
	require.Nil(t, err, "weekly cron request failed")
	require.Equal(t, http.StatusOK, weeklyCronResponse.StatusCode)
	defer weeklyCronResponse.Body.Close()

	var tasks []models.Task
	err = json.NewDecoder(weeklyCronResponse.Body).Decode(&tasks)
	require.Nil(t, err, "weekly cron endpoint returned unknown format", err)

	require.Equal(t, 1, len(tasks))

	ids := []*int{}

	for _, task := range tasks {
		ids = append(ids, &task.ID)
	}

	expectedTasks := expectedWeeklyTasks(ids)

	assert.Equal(t, len(tasks), len(expectedTasks))

	for i, task := range tasks {
		assert.Equal(t, expectedTasks[i].TemplateID, task.TemplateID)
		assert.Equal(t, expectedTasks[i].Status, task.Status)
		assert.Equal(t, expectedTasks[i].Notes, task.Notes)
		assert.Equal(t, expectedTasks[i].StartedAt, task.StartedAt)
		assert.Equal(t, expectedTasks[i].EndedAt, task.EndedAt)
		assert.Equal(t, expectedTasks[i].Date, task.Date)
	}

	getTasksRequest, err := http.NewRequest(http.MethodGet, tasksCollectionURL(), nil)
	require.Nil(t, err, "cannot create get tasks request", err)

	getTasksResponse, err := testClient.Do(getTasksRequest)
	require.Nil(t, err, "get request failed")
	require.Equal(t, http.StatusOK, getTasksResponse.StatusCode)
	defer getTasksResponse.Body.Close()

	var getTasksResult models.TaskListResponse
	err = json.NewDecoder(getTasksResponse.Body).Decode(&getTasksResult)
	require.Nil(t, err, "get tasks endpoint returned unknown format", err)

	assert.Equal(t, tasks, getTasksResult.Data)
}

func TestScheduleMonthlyCronHook(t *testing.T) {
	TruncateDB(t)

	bodies := validMonthlyTemplates()

	for _, body := range bodies {
		t.Run(body.Name, func(t *testing.T) {
			templateJson, err := json.Marshal(body)
			require.Nil(t, err, "cannot marshal the body", err)

			request, err := http.NewRequest(http.MethodPost, templatesCollectionURL(), bytes.NewReader(templateJson))
			require.Nil(t, err, "cannot create create request", err)

			response, err := testClient.Do(request)
			require.Nil(t, err, "create request failed", err)

			defer response.Body.Close()

			require.Equal(t, http.StatusCreated, response.StatusCode)
		})
	}

	getCollectionRequest, err := http.NewRequest(http.MethodGet, templatesCollectionURL(), nil)
	require.Nil(t, err, "cannot create get template request", err)

	getCollectionResponse, err := testClient.Do(getCollectionRequest)
	require.Nil(t, err, "get request failed")
	require.Equal(t, http.StatusOK, getCollectionResponse.StatusCode)
	defer getCollectionResponse.Body.Close()

	var templates models.TemplateListResponse
	err = json.NewDecoder(getCollectionResponse.Body).Decode(&templates)
	require.Nil(t, err, "get templates endpoint returned unknown format", err)

	assert.Equal(t, len(bodies), len(templates.Data))

	for i, template := range templates.Data {
		assert.Equal(t, template.Name, bodies[i].Name)
		assert.Equal(t, template.Description, bodies[i].Description)
		assert.Equal(t, template.StartsAt, bodies[i].StartsAt)
		assert.Equal(t, template.EndsAt, bodies[i].EndsAt)
		assert.Equal(t, template.Scheduling, bodies[i].Scheduling)
	}

	monthlyCronRequest, err := http.NewRequest(http.MethodPost, monthlyCronURL(), nil)
	require.Nil(t, err, "cannot create monthly cron request", err)

	monthlyCronResponse, err := testClient.Do(monthlyCronRequest)
	require.Nil(t, err, "monthly cron request failed")
	require.Equal(t, http.StatusOK, monthlyCronResponse.StatusCode)
	defer monthlyCronResponse.Body.Close()

	var tasks []models.Task
	err = json.NewDecoder(monthlyCronResponse.Body).Decode(&tasks)
	require.Nil(t, err, "monthly cron endpoint returned unknown format", err)

	require.Equal(t, 1, len(tasks))

	ids := []*int{}

	for _, task := range tasks {
		ids = append(ids, &task.ID)
	}

	expectedTasks := expectedMonthlyTasks(ids)

	assert.Equal(t, len(tasks), len(expectedTasks))

	for i, task := range tasks {
		assert.Equal(t, expectedTasks[i].TemplateID, task.TemplateID)
		assert.Equal(t, expectedTasks[i].Status, task.Status)
		assert.Equal(t, expectedTasks[i].Notes, task.Notes)
		assert.Equal(t, expectedTasks[i].StartedAt, task.StartedAt)
		assert.Equal(t, expectedTasks[i].EndedAt, task.EndedAt)
		assert.Equal(t, *expectedTasks[i].Date, *task.Date)
	}

	getTasksRequest, err := http.NewRequest(http.MethodGet, tasksCollectionURL(), nil)
	require.Nil(t, err, "cannot create get tasks request", err)

	getTasksResponse, err := testClient.Do(getTasksRequest)
	require.Nil(t, err, "get request failed")
	require.Equal(t, http.StatusOK, getTasksResponse.StatusCode)
	defer getTasksResponse.Body.Close()

	var getTasksResult models.TaskListResponse
	err = json.NewDecoder(getTasksResponse.Body).Decode(&getTasksResult)
	require.Nil(t, err, "get tasks endpoint returned unknown format", err)

	assert.Equal(t, tasks, getTasksResult.Data)
}
