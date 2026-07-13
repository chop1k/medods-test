package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/chop1k/medods-test/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func tasksCollectionURL() string {
	return testURL + "/v1/tasks/tasks"
}

func taskURL(id int) string {
	return testURL + "/v1/tasks/tasks/" + strconv.Itoa(id)
}

func TestTaskFinished(t *testing.T) {
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

	task := tasks[0]

	startedAt := time.Now().Format(time.TimeOnly)

	updateTaskToRunningBody := map[string]any{
		"status":     models.TaskStatusRunning,
		"started_at": startedAt,
	}
	updateTaskToRunningJson, err := json.Marshal(updateTaskToRunningBody)
	require.Nil(t, err, "cannot marshal the body", err)

	updateTaskToRunningRequest, err := http.NewRequest(http.MethodPatch, taskURL(task.ID), bytes.NewReader(updateTaskToRunningJson))
	require.Nil(t, err, "cannot create get tasks request", err)

	updateTaskToRunningResponse, err := testClient.Do(updateTaskToRunningRequest)
	require.Nil(t, err, "update task to running request failed")
	require.Equal(t, http.StatusOK, updateTaskToRunningResponse.StatusCode)
	defer updateTaskToRunningResponse.Body.Close()

	var updatedRunningTask models.Task
	err = json.NewDecoder(updateTaskToRunningResponse.Body).Decode(&updatedRunningTask)
	require.Nil(t, err, "update task endpoint returned unknown format", err)

	expectedRunningTask := models.Task{
		ID: *ids[0],
		TaskBody: models.TaskBody{
			TemplateID: task.TemplateID,
			MovedId:    nil,
			Status:     models.TaskStatusRunning,
			Notes:      nil,
			Date:       task.Date,
			StartedAt:  &startedAt,
			EndedAt:    task.EndedAt,
		},
	}

	assert.Equal(t, expectedRunningTask.ID, updatedRunningTask.ID)
	assert.Equal(t, expectedRunningTask.TemplateID, updatedRunningTask.TemplateID)
	assert.Equal(t, expectedRunningTask.Status, updatedRunningTask.Status)
	assert.Equal(t, expectedRunningTask.Notes, updatedRunningTask.Notes)
	assert.Equal(t, *expectedRunningTask.StartedAt, *updatedRunningTask.StartedAt)
	assert.Equal(t, expectedRunningTask.EndedAt, updatedRunningTask.EndedAt)
	assert.Equal(t, *expectedRunningTask.Date, *updatedRunningTask.Date)

	endedAt := time.Now().Add(time.Hour).Format(time.TimeOnly)

	updateTaskToFinishedBody := map[string]any{
		"status":   models.TaskStatusFinished,
		"ended_at": endedAt,
	}
	updateTaskToFinishedJson, err := json.Marshal(updateTaskToFinishedBody)
	require.Nil(t, err, "cannot marshal the body", err)

	updateTaskToFinishedRequest, err := http.NewRequest(http.MethodPatch, taskURL(task.ID), bytes.NewReader(updateTaskToFinishedJson))
	require.Nil(t, err, "cannot create get tasks request", err)

	updateTaskToFinishedResponse, err := testClient.Do(updateTaskToFinishedRequest)
	require.Nil(t, err, "update task to running request failed")
	require.Equal(t, http.StatusOK, updateTaskToFinishedResponse.StatusCode)
	defer updateTaskToFinishedResponse.Body.Close()

	var updatedFinishedTask models.Task
	err = json.NewDecoder(updateTaskToFinishedResponse.Body).Decode(&updatedFinishedTask)
	require.Nil(t, err, "update task endpoint returned unknown format", err)

	expectedFinishedTask := models.Task{
		ID: *ids[0],
		TaskBody: models.TaskBody{
			TemplateID: task.TemplateID,
			MovedId:    nil,
			Status:     models.TaskStatusFinished,
			Notes:      nil,
			Date:       task.Date,
			StartedAt:  &startedAt,
			EndedAt:    &endedAt,
		},
	}

	assert.Equal(t, expectedFinishedTask.ID, updatedFinishedTask.ID)
	assert.Equal(t, expectedFinishedTask.TemplateID, updatedFinishedTask.TemplateID)
	assert.Equal(t, expectedFinishedTask.Status, updatedFinishedTask.Status)
	assert.Equal(t, expectedFinishedTask.Notes, updatedFinishedTask.Notes)
	assert.Equal(t, expectedFinishedTask.StartedAt, updatedFinishedTask.StartedAt)
	assert.Equal(t, expectedFinishedTask.EndedAt, updatedFinishedTask.EndedAt)
	assert.Equal(t, expectedFinishedTask.Date, updatedFinishedTask.Date)
}

func TestTaskPendingCancelled(t *testing.T) {
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

	task := tasks[0]

	updateTaskToCancelledBody := map[string]any{
		"status": models.TaskStatusCancelled,
	}
	updateTaskToCancelledJson, err := json.Marshal(updateTaskToCancelledBody)
	require.Nil(t, err, "cannot marshal the body", err)

	updateTaskToCancelledRequest, err := http.NewRequest(http.MethodPatch, taskURL(task.ID), bytes.NewReader(updateTaskToCancelledJson))
	require.Nil(t, err, "cannot create get tasks request", err)

	updateTaskToCancelledResponse, err := testClient.Do(updateTaskToCancelledRequest)
	require.Nil(t, err, "update task to running request failed")
	require.Equal(t, http.StatusOK, updateTaskToCancelledResponse.StatusCode)
	defer updateTaskToCancelledResponse.Body.Close()

	var updatedCancelledTask models.Task
	err = json.NewDecoder(updateTaskToCancelledResponse.Body).Decode(&updatedCancelledTask)
	require.Nil(t, err, "update task endpoint returned unknown format", err)

	expectedRunningTask := models.Task{
		ID: *ids[0],
		TaskBody: models.TaskBody{
			TemplateID: task.TemplateID,
			MovedId:    nil,
			Status:     models.TaskStatusCancelled,
			Notes:      nil,
			Date:       task.Date,
			StartedAt:  nil,
			EndedAt:    nil,
		},
	}

	assert.Equal(t, expectedRunningTask.ID, updatedCancelledTask.ID)
	assert.Equal(t, expectedRunningTask.TemplateID, updatedCancelledTask.TemplateID)
	assert.Equal(t, expectedRunningTask.Status, updatedCancelledTask.Status)
	assert.Equal(t, expectedRunningTask.Notes, updatedCancelledTask.Notes)
	assert.Equal(t, expectedRunningTask.StartedAt, updatedCancelledTask.StartedAt)
	assert.Equal(t, expectedRunningTask.EndedAt, updatedCancelledTask.EndedAt)
	assert.Equal(t, expectedRunningTask.Date, updatedCancelledTask.Date)
}

func TestTaskRunningCancelled(t *testing.T) {
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

	task := tasks[0]

	startedAt := time.Now().Format(time.TimeOnly)

	updateTaskToRunningBody := map[string]any{
		"status":     models.TaskStatusRunning,
		"started_at": startedAt,
	}
	updateTaskToRunningJson, err := json.Marshal(updateTaskToRunningBody)
	require.Nil(t, err, "cannot marshal the body", err)

	updateTaskToRunningRequest, err := http.NewRequest(http.MethodPatch, taskURL(task.ID), bytes.NewReader(updateTaskToRunningJson))
	require.Nil(t, err, "cannot create get tasks request", err)

	updateTaskToRunningResponse, err := testClient.Do(updateTaskToRunningRequest)
	require.Nil(t, err, "update task to running request failed")
	require.Equal(t, http.StatusOK, updateTaskToRunningResponse.StatusCode)
	defer updateTaskToRunningResponse.Body.Close()

	var updatedRunningTask models.Task
	err = json.NewDecoder(updateTaskToRunningResponse.Body).Decode(&updatedRunningTask)
	require.Nil(t, err, "update task endpoint returned unknown format", err)

	expectedRunningTask := models.Task{
		ID: *ids[0],
		TaskBody: models.TaskBody{
			TemplateID: task.TemplateID,
			MovedId:    nil,
			Status:     models.TaskStatusRunning,
			Notes:      nil,
			Date:       task.Date,
			StartedAt:  &startedAt,
			EndedAt:    task.EndedAt,
		},
	}

	assert.Equal(t, expectedRunningTask.ID, updatedRunningTask.ID)
	assert.Equal(t, expectedRunningTask.TemplateID, updatedRunningTask.TemplateID)
	assert.Equal(t, expectedRunningTask.Status, updatedRunningTask.Status)
	assert.Equal(t, expectedRunningTask.Notes, updatedRunningTask.Notes)
	assert.Equal(t, expectedRunningTask.StartedAt, updatedRunningTask.StartedAt)
	assert.Equal(t, expectedRunningTask.EndedAt, updatedRunningTask.EndedAt)
	assert.Equal(t, *expectedRunningTask.Date, *updatedRunningTask.Date)

	endedAt := time.Now().Add(time.Hour).Format(time.TimeOnly)

	updateTaskToCancelledBody := map[string]any{
		"status":   models.TaskStatusCancelled,
		"ended_at": endedAt,
	}
	updateTaskToCancelledJson, err := json.Marshal(updateTaskToCancelledBody)
	require.Nil(t, err, "cannot marshal the body", err)

	updateTaskToCancelledRequest, err := http.NewRequest(http.MethodPatch, taskURL(task.ID), bytes.NewReader(updateTaskToCancelledJson))
	require.Nil(t, err, "cannot create get tasks request", err)

	updateTaskToCancelledResponse, err := testClient.Do(updateTaskToCancelledRequest)
	require.Nil(t, err, "update task to running request failed")
	require.Equal(t, http.StatusOK, updateTaskToCancelledResponse.StatusCode)
	defer updateTaskToCancelledResponse.Body.Close()

	var updatedCancelledTask models.Task
	err = json.NewDecoder(updateTaskToCancelledResponse.Body).Decode(&updatedCancelledTask)
	require.Nil(t, err, "update task endpoint returned unknown format", err)

	expectedCancelledTask := models.Task{
		ID: *ids[0],
		TaskBody: models.TaskBody{
			TemplateID: task.TemplateID,
			MovedId:    nil,
			Status:     models.TaskStatusCancelled,
			Notes:      nil,
			Date:       task.Date,
			StartedAt:  &startedAt,
			EndedAt:    &endedAt,
		},
	}

	assert.Equal(t, expectedCancelledTask.ID, updatedCancelledTask.ID)
	assert.Equal(t, expectedCancelledTask.TemplateID, updatedCancelledTask.TemplateID)
	assert.Equal(t, expectedCancelledTask.Status, updatedCancelledTask.Status)
	assert.Equal(t, expectedCancelledTask.Notes, updatedCancelledTask.Notes)
	assert.Equal(t, expectedCancelledTask.StartedAt, updatedCancelledTask.StartedAt)
	assert.Equal(t, expectedCancelledTask.EndedAt, updatedCancelledTask.EndedAt)
	assert.Equal(t, expectedCancelledTask.Date, updatedCancelledTask.Date)
}

func TestTaskPendingMoved(t *testing.T) {
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

	task := tasks[0]

	now := time.Now()
	oldDate := now.Format("02-01-2006")
	date := now.AddDate(0, 0, 1).Format("02-01-2006")

	updateTaskToMovedBody := map[string]any{
		"status": models.TaskStatusMoved,
		"date":   date,
	}
	updateTaskToMovedJson, err := json.Marshal(updateTaskToMovedBody)
	require.Nil(t, err, "cannot marshal the body", err)

	updateTaskToMovedRequest, err := http.NewRequest(http.MethodPatch, taskURL(task.ID), bytes.NewReader(updateTaskToMovedJson))
	require.Nil(t, err, "cannot create get tasks request", err)

	updateTaskToMovedResponse, err := testClient.Do(updateTaskToMovedRequest)
	require.Nil(t, err, "update task to running request failed")
	require.Equal(t, http.StatusOK, updateTaskToMovedResponse.StatusCode)
	defer updateTaskToMovedResponse.Body.Close()

	var updatedMovedTask models.Task
	err = json.NewDecoder(updateTaskToMovedResponse.Body).Decode(&updatedMovedTask)
	require.Nil(t, err, "update task endpoint returned unknown format", err)

	expectedMovedTask := models.Task{
		ID: *ids[0],
		TaskBody: models.TaskBody{
			TemplateID: task.TemplateID,
			MovedId:    nil,
			Status:     models.TaskStatusMoved,
			Notes:      nil,
			Date:       &oldDate,
			StartedAt:  nil,
			EndedAt:    nil,
		},
	}

	assert.Equal(t, expectedMovedTask.ID, updatedMovedTask.ID)
	assert.Equal(t, expectedMovedTask.TemplateID, updatedMovedTask.TemplateID)
	assert.Equal(t, expectedMovedTask.Status, updatedMovedTask.Status)
	assert.Equal(t, expectedMovedTask.Notes, updatedMovedTask.Notes)
	assert.Equal(t, expectedMovedTask.StartedAt, updatedMovedTask.StartedAt)
	assert.Equal(t, expectedMovedTask.EndedAt, updatedMovedTask.EndedAt)
	assert.Equal(t, *expectedMovedTask.Date, *updatedMovedTask.Date)
}
