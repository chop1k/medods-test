package e2e

// import (
// 	"bytes"
// 	"encoding/json"
// 	"math/rand"
// 	"net/http"
// 	"strconv"
// 	"testing"

// 	"github.com/chop1k/medods-test/internal/models"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// const dailyCronURL = "/v1/scheduling/daily-cron-url"

// func tasksCollectionURL() string {
// 	return testURL + "/v1/tasks/tasks"
// }

// func taskURL(id int) string {
// 	return testURL + "/v1/tasks/tasks/" + strconv.Itoa(id)
// }

// func validTaskBodies() []models.TaskBody {
// 	return []models.TaskBody{}
// }

// func TestGetTasks(t *testing.T) {
// 	TruncateDB(t)

// 	bodies := validTaskBodies()

// 	for _, body := range bodies {
// 		t.Run(body.Name, func(t *testing.T) {
// 			templateJson, err := json.Marshal(body)
// 			require.Nil(t, err, "cannot marshal the body", err)

// 			request, err := http.NewRequest(http.MethodPost, tasksCollectionURL(), bytes.NewReader(templateJson))
// 			require.Nil(t, err, "cannot create create request", err)

// 			response, err := testClient.Do(request)
// 			require.Nil(t, err, "create request failed", err)

// 			defer response.Body.Close()

// 			require.Equal(t, http.StatusCreated, response.StatusCode)
// 		})
// 	}

// 	request, err := http.NewRequest(http.MethodGet, tasksCollectionURL(), nil)
// 	require.Nil(t, err, "cannot create get template request", err)

// 	response, err := testClient.Do(request)
// 	require.Nil(t, err, "get request failed")
// 	require.Equal(t, http.StatusOK, response.StatusCode)
// 	defer response.Body.Close()

// 	var templates models.TemplateListResponse
// 	err = json.NewDecoder(response.Body).Decode(&templates)
// 	require.Nil(t, err, "get templates endpoint returned unknown format", err)

// 	assert.Equal(t, len(bodies), len(templates.Data))

// 	for i, template := range templates.Data {
// 		assert.Equal(t, template.Name, bodies[i].Name)
// 		assert.Equal(t, template.Description, bodies[i].Description)
// 		assert.Equal(t, template.StartsAt, bodies[i].StartsAt)
// 		assert.Equal(t, template.EndsAt, bodies[i].EndsAt)
// 		assert.Equal(t, template.Scheduling, bodies[i].Scheduling)
// 	}
// }

// func TestUpdateTasks(t *testing.T) {
// 	TruncateDB(t)

// 	bodies := validTaskBodies()

// 	for _, body := range bodies {
// 		t.Run(body.Name, func(t *testing.T) {
// 			templateJson, err := json.Marshal(body)
// 			require.Nil(t, err, "cannot marshal the body", err)

// 			request, err := http.NewRequest(http.MethodPost, tasksCollectionURL(), bytes.NewReader(templateJson))
// 			require.Nil(t, err, "cannot create create request", err)

// 			response, err := testClient.Do(request)
// 			require.Nil(t, err, "create request failed", err)

// 			defer response.Body.Close()

// 			require.Equal(t, http.StatusCreated, response.StatusCode)
// 		})
// 	}

// 	getCollectionRequest, err := http.NewRequest(http.MethodGet, tasksCollectionURL(), nil)
// 	require.Nil(t, err, "cannot create get template request", err)

// 	getCollectionResponse, err := testClient.Do(getCollectionRequest)
// 	require.Nil(t, err, "get request failed")
// 	require.Equal(t, http.StatusOK, getCollectionResponse.StatusCode)
// 	defer getCollectionResponse.Body.Close()

// 	var templates models.TemplateListResponse
// 	err = json.NewDecoder(getCollectionResponse.Body).Decode(&templates)
// 	require.Nil(t, err, "get templates endpoint returned unknown format", err)

// 	assert.Equal(t, len(bodies), len(templates.Data))

// 	for i, template := range templates.Data {
// 		assert.Equal(t, template.Name, bodies[i].Name)
// 		assert.Equal(t, template.Description, bodies[i].Description)
// 		assert.Equal(t, template.StartsAt, bodies[i].StartsAt)
// 		assert.Equal(t, template.EndsAt, bodies[i].EndsAt)
// 		assert.Equal(t, template.Enabled, bodies[i].Enabled)
// 		assert.Equal(t, template.Scheduling, bodies[i].Scheduling)
// 	}

// 	templateNumber := rand.New(rand.NewSource(testSeed)).Intn(len(templates.Data))
// 	templateID := templates.Data[templateNumber].ID

// 	validUpdateBody := map[string]any{
// 		"enabled": false,
// 	}
// 	validUpdateJson, err := json.Marshal(validUpdateBody)
// 	require.Nil(t, err, "cannot marshal the update valid body", err)

// 	validUpdateRequest, err := http.NewRequest(http.MethodPut, taskURL(templateID), bytes.NewReader(validUpdateJson))
// 	validUpdateRequest.Header.Add("content-type", "application/json")
// 	require.Nil(t, err, "cannot create valid update template request", err)

// 	validUpdateResponse, err := testClient.Do(validUpdateRequest)
// 	require.Nil(t, err, "valid update request failed")
// 	require.Equal(t, http.StatusOK, validUpdateResponse.StatusCode)
// 	defer validUpdateResponse.Body.Close()

// 	invalidUpdateBody := map[string]any{
// 		"name": "321",
// 	}
// 	invalidUpdateJson, err := json.Marshal(invalidUpdateBody)
// 	require.Nil(t, err, "cannot marshal the update invalid body", err)

// 	invalidUpdateRequest, err := http.NewRequest(http.MethodPut, taskURL(templateID), bytes.NewReader(invalidUpdateJson))
// 	invalidUpdateRequest.Header.Add("content-type", "application/json")
// 	require.Nil(t, err, "cannot create invalud update template request", err)

// 	invalidUpdateResponse, err := testClient.Do(invalidUpdateRequest)
// 	require.Nil(t, err, "valid update request failed")
// 	require.Equal(t, http.StatusBadRequest, invalidUpdateResponse.StatusCode)
// 	defer invalidUpdateResponse.Body.Close()

// 	getTemplateRequest, err := http.NewRequest(http.MethodGet, taskURL(templateID), nil)
// 	require.Nil(t, err, "cannot create get template request", err)

// 	getTemplateResponse, err := testClient.Do(getTemplateRequest)
// 	require.Nil(t, err, "get request failed")
// 	require.Equal(t, http.StatusOK, getTemplateResponse.StatusCode)
// 	defer getTemplateResponse.Body.Close()

// 	var template models.Template
// 	err = json.NewDecoder(getTemplateResponse.Body).Decode(&template)
// 	require.Nil(t, err, "get template endpoint returned unknown format", err)

// 	assert.Equal(t, template.Name, bodies[templateNumber].Name)
// 	assert.Equal(t, template.Description, bodies[templateNumber].Description)
// 	assert.Equal(t, template.StartsAt, bodies[templateNumber].StartsAt)
// 	assert.Equal(t, template.EndsAt, bodies[templateNumber].EndsAt)
// 	assert.Equal(t, template.Enabled, false)
// 	assert.Equal(t, template.Scheduling, bodies[templateNumber].Scheduling)
// }

// func TestRemoveTasks(t *testing.T) {
// 	TruncateDB(t)

// 	bodies := validTaskBodies()

// 	for _, body := range bodies {
// 		t.Run(body.Name, func(t *testing.T) {
// 			templateJson, err := json.Marshal(body)
// 			require.Nil(t, err, "cannot marshal the body", err)

// 			request, err := http.NewRequest(http.MethodPost, tasksCollectionURL(), bytes.NewReader(templateJson))
// 			require.Nil(t, err, "cannot create create request", err)

// 			response, err := testClient.Do(request)
// 			require.Nil(t, err, "create request failed", err)

// 			defer response.Body.Close()

// 			require.Equal(t, http.StatusCreated, response.StatusCode)
// 		})
// 	}

// 	getCollectionRequest, err := http.NewRequest(http.MethodGet, tasksCollectionURL(), nil)
// 	require.Nil(t, err, "cannot create get template request", err)

// 	getCollectionResponse, err := testClient.Do(getCollectionRequest)
// 	require.Nil(t, err, "get request failed")
// 	require.Equal(t, http.StatusOK, getCollectionResponse.StatusCode)
// 	defer getCollectionResponse.Body.Close()

// 	var templates models.TemplateListResponse
// 	err = json.NewDecoder(getCollectionResponse.Body).Decode(&templates)
// 	require.Nil(t, err, "get templates endpoint returned unknown format", err)

// 	assert.Equal(t, len(bodies), len(templates.Data))

// 	for i, template := range templates.Data {
// 		assert.Equal(t, template.Name, bodies[i].Name)
// 		assert.Equal(t, template.Description, bodies[i].Description)
// 		assert.Equal(t, template.StartsAt, bodies[i].StartsAt)
// 		assert.Equal(t, template.EndsAt, bodies[i].EndsAt)
// 		assert.Equal(t, template.Scheduling, bodies[i].Scheduling)
// 	}

// 	templateNumber := rand.New(rand.NewSource(testSeed)).Intn(len(templates.Data))
// 	templateID := templates.Data[templateNumber].ID

// 	removeRequest, err := http.NewRequest(http.MethodDelete, taskURL(templateID), nil)
// 	require.Nil(t, err, "cannot create delete template request", err)

// 	removeResponse, err := testClient.Do(removeRequest)
// 	require.Nil(t, err, "remove request failed")
// 	require.Equal(t, http.StatusNoContent, removeResponse.StatusCode)
// 	defer removeResponse.Body.Close()

// 	getTemplateRequest, err := http.NewRequest(http.MethodGet, taskURL(templateID), nil)
// 	require.Nil(t, err, "cannot create get template request", err)

// 	getTemplateResponse, err := testClient.Do(getTemplateRequest)
// 	require.Nil(t, err, "get request failed")
// 	require.Equal(t, http.StatusNotFound, getTemplateResponse.StatusCode)
// 	defer getTemplateResponse.Body.Close()
// }
