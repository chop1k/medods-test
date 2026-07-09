package e2e

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/chop1k/medods-test/internal/models"
)

func validTagBody() models.TagBody {
	return models.TagBody{
		Name: "Household",
		Type: models.TagTypeUserDefined,
	}
}

func TestGetTags_Success(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	var out models.TagListResponse
	resp := doJSON(t, client, http.MethodGet, srv.URL+"/v1/tasks/tags", nil, &out)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, out.Data)
	assert.Equal(t, 1, out.Meta.Page)
	assert.Equal(t, 20, out.Meta.Limit)
}

func TestGetTags_WithPaginationAndSorting(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	url := srv.URL + "/v1/tasks/tags?page=2&limit=10&sort=desc&sort-field=name"
	var out models.TagListResponse
	resp := doJSON(t, client, http.MethodGet, url, nil, &out)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, out.Meta.Page)
	assert.Equal(t, 10, out.Meta.Limit)
}

func TestGetTags_InvalidQuery(t *testing.T) {
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
			resp := doJSON(t, client, http.MethodGet, srv.URL+"/v1/tasks/tags"+qs, nil, &out)

			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
			assert.Equal(t, http.StatusBadRequest, out.Status)
		})
	}
}

func TestCreateTag_Success(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	body := validTagBody()

	var out models.Tag
	resp := doJSON(t, client, http.MethodPost, srv.URL+"/v1/tasks/tags", body, &out)

	require.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, body.Name, out.Name)
	assert.Equal(t, body.Type, out.Type)
}

func TestCreateTag_ValidationError(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	cases := map[string]models.TagBody{
		"name too short": {
			Name: "a",
			Type: models.TagTypePredefined,
		},
		"missing type": {
			Name: "Valid name",
		},
		"invalid type": {
			Name: "Valid name",
			Type: "not-a-type",
		},
	}

	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			var out models.ValidationErrorResponse
			resp := doJSON(t, client, http.MethodPost, srv.URL+"/v1/tasks/tags", body, &out)

			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
			assert.Equal(t, http.StatusBadRequest, out.Status)
		})
	}
}

func TestGetTagByID_Success(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	var out models.Tag
	resp := doJSON(t, client, http.MethodGet, srv.URL+"/v1/tasks/tags/42", nil, &out)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, int64(42), out.ID)
}

func TestGetTagByID_InvalidID(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	cases := []string{"0", "-1", "not-a-number"}

	for _, id := range cases {
		t.Run(id, func(t *testing.T) {
			var out models.ValidationErrorResponse
			resp := doJSON(t, client, http.MethodGet, fmt.Sprintf("%s/v1/tasks/tags/%s", srv.URL, id), nil, &out)

			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
			assert.Equal(t, http.StatusBadRequest, out.Status)
		})
	}
}

func TestUpdateTag_Success(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	body := validTagBody()
	body.Name = "Updated name"

	var out models.Tag
	resp := doJSON(t, client, http.MethodPut, srv.URL+"/v1/tasks/tags/42", body, &out)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, int64(42), out.ID)
	assert.Equal(t, "Updated name", out.Name)
}

func TestUpdateTag_InvalidID(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	body := validTagBody()

	var out models.ValidationErrorResponse
	resp := doJSON(t, client, http.MethodPut, srv.URL+"/v1/tasks/tags/0", body, &out)

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, http.StatusBadRequest, out.Status)
}

func TestUpdateTag_ValidationError(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	invalidBody := map[string]any{
		"name": "a",
		"type": "predefined",
	}

	var out models.ValidationErrorResponse
	resp := doJSON(t, client, http.MethodPut, srv.URL+"/v1/tasks/tags/42", invalidBody, &out)

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, http.StatusBadRequest, out.Status)
}

func TestDeleteTag_Success(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	resp := doJSON(t, client, http.MethodDelete, srv.URL+"/v1/tasks/tags/42", nil, nil)

	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestDeleteTag_InvalidID(t *testing.T) {
	srv := newTestServer(t)
	client := srv.Client()

	var out models.ValidationErrorResponse
	resp := doJSON(t, client, http.MethodDelete, srv.URL+"/v1/tasks/tags/abc", nil, &out)

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, http.StatusBadRequest, out.Status)
}
