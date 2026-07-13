package e2e

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"testing"

	"github.com/chop1k/medods-test/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func tagsCollectionURL() string {
	return testURL + "/v1/grouping/tags"
}

func tagURL(id int) string {
	return testURL + "/v1/grouping/tags/" + strconv.Itoa(id)
}

func validTagBodies() []models.TagBody {
	description := "321"

	return []models.TagBody{
		{
			Name:        "123",
			Description: nil,
		},
		{
			Name:        "123",
			Description: &description,
		},
	}
}

func TestGetTags(t *testing.T) {
	TruncateDB(t)

	bodies := validTagBodies()

	for _, body := range bodies {
		t.Run(body.Name, func(t *testing.T) {
			tagJson, err := json.Marshal(body)
			require.Nil(t, err, "cannot marshal the body", err)

			request, err := http.NewRequest(http.MethodPost, tagsCollectionURL(), bytes.NewReader(tagJson))
			require.Nil(t, err, "cannot create create request", err)

			response, err := testClient.Do(request)
			require.Nil(t, err, "create request failed", err)
			defer response.Body.Close()

			require.Equal(t, http.StatusCreated, response.StatusCode)
		})
	}

	request, err := http.NewRequest(http.MethodGet, tagsCollectionURL(), nil)
	require.Nil(t, err, "cannot create get tags request", err)

	response, err := testClient.Do(request)
	require.Nil(t, err, "get request failed")
	require.Equal(t, http.StatusOK, response.StatusCode)

	var tags models.TagListResponse
	err = json.NewDecoder(response.Body).Decode(&tags)
	require.Nil(t, err, "get tags endpoint returned unknown format", err)

	assert.Equal(t, len(bodies), len(tags.Data))

	for i, tag := range tags.Data {
		assert.Equal(t, tag.Name, bodies[i].Name)
		assert.Equal(t, tag.Description, bodies[i].Description)
	}
}

func TestRemoveTags(t *testing.T) {
	TruncateDB(t)

	bodies := validTagBodies()

	for _, body := range bodies {
		t.Run(body.Name, func(t *testing.T) {
			tagJson, err := json.Marshal(body)
			require.Nil(t, err, "cannot marshal the body", err)

			request, err := http.NewRequest(http.MethodPost, tagsCollectionURL(), bytes.NewReader(tagJson))
			require.Nil(t, err, "cannot create create request", err)

			response, err := testClient.Do(request)
			require.Nil(t, err, "create request failed", err)

			defer response.Body.Close()

			require.Equal(t, http.StatusCreated, response.StatusCode)
		})
	}

	getCollectionRequest, err := http.NewRequest(http.MethodGet, tagsCollectionURL(), nil)
	require.Nil(t, err, "cannot create get tags request", err)

	getCollectionResponse, err := testClient.Do(getCollectionRequest)
	require.Nil(t, err, "get request failed")
	require.Equal(t, http.StatusOK, getCollectionResponse.StatusCode)
	defer getCollectionResponse.Body.Close()

	var tags models.TemplateListResponse
	err = json.NewDecoder(getCollectionResponse.Body).Decode(&tags)
	require.Nil(t, err, "get tags endpoint returned unknown format", err)

	assert.Equal(t, len(bodies), len(tags.Data))

	for i, template := range tags.Data {
		assert.Equal(t, template.Name, bodies[i].Name)
		assert.Equal(t, template.Description, bodies[i].Description)
	}

	tagNumber := rand.New(rand.NewSource(testSeed)).Intn(len(tags.Data))
	tagID := tags.Data[tagNumber].ID

	removeRequest, err := http.NewRequest(http.MethodDelete, tagURL(tagID), nil)
	require.Nil(t, err, "cannot create delete tag request", err)

	removeResponse, err := testClient.Do(removeRequest)
	require.Nil(t, err, "remove request failed")
	require.Equal(t, http.StatusNoContent, removeResponse.StatusCode)
	defer removeResponse.Body.Close()

	getTemplateRequest, err := http.NewRequest(http.MethodGet, tagURL(tagID), nil)
	require.Nil(t, err, "cannot create get tag request", err)

	getTemplateResponse, err := testClient.Do(getTemplateRequest)
	require.Nil(t, err, "get request failed")
	require.Equal(t, http.StatusNotFound, getTemplateResponse.StatusCode)
	defer removeResponse.Body.Close()
}
