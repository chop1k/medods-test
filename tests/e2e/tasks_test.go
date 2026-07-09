package e2e

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/chop1k/medods-test/internal/models"
)

func validTaskBody() models.TaskBody {
	return models.TaskBody{
		TemplateID: 7,
		Status:     models.TaskStatusPending,
		Notes:      "Some notes",
	}
}

func TestGetTasks_Success(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	var out models.TaskListResponse
	resp := doJSON(t, client, http.MethodGet, srv.URL+"/v1/tasks", nil, &out)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, out.Data)
	assert.Equal(t, 1, out.Meta.Page)
	assert.Equal(t, 20, out.Meta.Limit)
}

func TestGetTasks_WithPaginationAndSorting(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	url := srv.URL + "/v1/tasks?page=2&limit=10&sort=desc&sort-field=status"
	var out models.TaskListResponse
	resp := doJSON(t, client, http.MethodGet, url, nil, &out)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, out.Meta.Page)
	assert.Equal(t, 10, out.Meta.Limit)
}

func TestGetTasks_InvalidQuery(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	cases := []string{
		"?page=0",
		"?limit=101",
		"?sort=sideways",
		"?sort-field=unknown",
	}

	for _, qs := range cases {
		t.Run(qs, func(t *testing.T) {
			var out models.ValidationErrorResponse
			resp := doJSON(t, client, http.MethodGet, srv.URL+"/v1/tasks"+qs, nil, &out)

			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
			assert.Equal(t, http.StatusBadRequest, out.Status)
		})
	}
}

func TestCreateTask_Success(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	body := validTaskBody()

	var out models.Task
	resp := doJSON(t, client, http.MethodPost, srv.URL+"/v1/tasks", body, &out)

	require.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, body.TemplateID, out.TemplateID)
	assert.Equal(t, body.Status, out.Status)
}

func TestCreateTask_ValidationError(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	cases := map[string]models.TaskBody{
		"missing template_id": {
			Status: models.TaskStatusPending,
		},
		"invalid status": {
			TemplateID: 1,
			Status:     "not-a-status",
		},
		"notes too short": {
			TemplateID: 1,
			Notes:      "a",
		},
	}

	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			var out models.ValidationErrorResponse
			resp := doJSON(t, client, http.MethodPost, srv.URL+"/v1/tasks", body, &out)

			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
			assert.Equal(t, http.StatusBadRequest, out.Status)
		})
	}
}

func TestGetTaskByID_Success(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	var out models.Task
	resp := doJSON(t, client, http.MethodGet, srv.URL+"/v1/tasks/42", nil, &out)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, int64(42), out.ID)
}

func TestGetTaskByID_InvalidID(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	cases := []string{"0", "-1", "not-a-number"}

	for _, id := range cases {
		t.Run(id, func(t *testing.T) {
			var out models.ValidationErrorResponse
			resp := doJSON(t, client, http.MethodGet, fmt.Sprintf("%s/v1/tasks/%s", srv.URL, id), nil, &out)

			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
			assert.Equal(t, http.StatusBadRequest, out.Status)
		})
	}
}

func TestUpdateTask_Success(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	body := validTaskBody()
	body.Status = models.TaskStatusRunning

	var out models.Task
	resp := doJSON(t, client, http.MethodPut, srv.URL+"/v1/tasks/42", body, &out)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, int64(42), out.ID)
	assert.Equal(t, models.TaskStatusRunning, out.Status)
}

func TestUpdateTask_InvalidID(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	body := validTaskBody()

	var out models.ValidationErrorResponse
	resp := doJSON(t, client, http.MethodPut, srv.URL+"/v1/tasks/0", body, &out)

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, http.StatusBadRequest, out.Status)
}

func TestUpdateTask_ValidationError(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	invalidBody := map[string]any{
		"status": "pending",
		"notes":  "valid notes",
	}

	var out models.ValidationErrorResponse
	resp := doJSON(t, client, http.MethodPut, srv.URL+"/v1/tasks/42", invalidBody, &out)

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, http.StatusBadRequest, out.Status)
}

func TestDeleteTask_Success(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	resp := doJSON(t, client, http.MethodDelete, srv.URL+"/v1/tasks/42", nil, nil)

	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestDeleteTask_InvalidID(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	var out models.ValidationErrorResponse
	resp := doJSON(t, client, http.MethodDelete, srv.URL+"/v1/tasks/abc", nil, &out)

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, http.StatusBadRequest, out.Status)
}
