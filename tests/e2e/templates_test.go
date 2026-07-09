package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/chop1k/medods-test/internal/models"
)

func validTemplateBody() models.TemplateBody {
	return models.TemplateBody{
		Name:        "Water the plants",
		Description: "Daily watering routine",
		StartsAt:    time.Now(),
		EndsAt:      time.Now().Add(24 * time.Hour),
	}
}

func TestGetTemplates_Success(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	var out models.TemplateListResponse
	resp := doJSON(t, client, http.MethodGet, srv.URL+"/v1/tasks/templates", nil, &out)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, out.Data)
	assert.Equal(t, 1, out.Meta.Page)
	assert.Equal(t, 20, out.Meta.Limit)
}

func TestGetTemplates_WithPaginationAndSorting(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	url := srv.URL + "/v1/tasks/templates?page=2&limit=10&sort=desc&sort-field=name"
	var out models.TemplateListResponse
	resp := doJSON(t, client, http.MethodGet, url, nil, &out)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, out.Meta.Page)
	assert.Equal(t, 10, out.Meta.Limit)
}

func TestGetTemplates_InvalidQuery(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	cases := []string{
		"?page=0",             // below minimum
		"?limit=101",          // above maximum
		"?sort=sideways",      // not in enum
		"?sort-field=unknown", // not in enum
	}

	for _, qs := range cases {
		t.Run(qs, func(t *testing.T) {
			var out models.ValidationErrorResponse
			resp := doJSON(t, client, http.MethodGet, srv.URL+"/v1/tasks/templates"+qs, nil, &out)

			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
			assert.Equal(t, http.StatusBadRequest, out.Status)
		})
	}
}

func TestCreateTemplate_Success(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	body := validTemplateBody()

	var out models.Template
	resp := doJSON(t, client, http.MethodPost, srv.URL+"/v1/tasks/templates", body, &out)

	require.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, body.Name, out.Name)
	assert.Equal(t, body.Description, out.Description)
}

func TestCreateTemplate_ValidationError(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	cases := map[string]models.TemplateBody{
		"name too short": {
			Name:        "a",
			Description: "valid description",
			StartsAt:    time.Now(),
			EndsAt:      time.Now().Add(time.Hour),
		},
		"description too short": {
			Name:        "Valid name",
			Description: "a",
			StartsAt:    time.Now(),
			EndsAt:      time.Now().Add(time.Hour),
		},
	}

	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			var out models.ValidationErrorResponse
			resp := doJSON(t, client, http.MethodPost, srv.URL+"/v1/tasks/templates", body, &out)

			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
			assert.Equal(t, http.StatusBadRequest, out.Status)
		})
	}
}

func TestCreateTemplate_MissingRequiredFields(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	var out models.ValidationErrorResponse
	resp := doJSON(t, client, http.MethodPost, srv.URL+"/v1/tasks/templates", map[string]any{}, &out)

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, http.StatusBadRequest, out.Status)
}

func TestGetTemplateByID_Success(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	var out models.Template
	resp := doJSON(t, client, http.MethodGet, srv.URL+"/v1/tasks/templates/42", nil, &out)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, int64(42), out.ID)
}

func TestGetTemplateByID_InvalidID(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	cases := []string{"0", "-1", "not-a-number"}

	for _, id := range cases {
		t.Run(id, func(t *testing.T) {
			var out models.ValidationErrorResponse
			resp := doJSON(t, client, http.MethodGet, fmt.Sprintf("%s/v1/tasks/templates/%s", srv.URL, id), nil, &out)

			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
			assert.Equal(t, http.StatusBadRequest, out.Status)
		})
	}
}

func TestUpdateTemplate_Success(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	body := validTemplateBody()
	body.Name = "Updated name"

	var out models.Template
	resp := doJSON(t, client, http.MethodPut, srv.URL+"/v1/tasks/templates/42", body, &out)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, int64(42), out.ID)
	assert.Equal(t, "Updated name", out.Name)
}

func TestUpdateTemplate_InvalidID(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	body := validTemplateBody()

	var out models.ValidationErrorResponse
	resp := doJSON(t, client, http.MethodPut, srv.URL+"/v1/tasks/templates/0", body, &out)

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, http.StatusBadRequest, out.Status)
}

func TestUpdateTemplate_ValidationError(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	invalidBody := map[string]any{
		"name":        "ok",
		"description": "a",
		"starts_at":   time.Now(),
		"ends_at":     time.Now().Add(time.Hour),
	}

	var out models.ValidationErrorResponse
	resp := doJSON(t, client, http.MethodPut, srv.URL+"/v1/tasks/templates/42", invalidBody, &out)

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, http.StatusBadRequest, out.Status)
}

func TestDeleteTemplate_Success(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	resp := doJSON(t, client, http.MethodDelete, srv.URL+"/v1/tasks/templates/42", nil, nil)

	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestDeleteTemplate_InvalidID(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	var out models.ValidationErrorResponse
	resp := doJSON(t, client, http.MethodDelete, srv.URL+"/v1/tasks/templates/abc", nil, &out)

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, http.StatusBadRequest, out.Status)
}
