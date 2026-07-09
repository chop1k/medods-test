package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/chop1k/medods-test/internal/models"
	"github.com/chop1k/medods-test/internal/responses"
)

type TagHandler struct {
}

func NewTagHandler() *TagHandler {
	return &TagHandler{}
}

func (h *TagHandler) GetTags(c *gin.Context) {
	var query models.ListTagsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		responses.ValidationError(c, err)
		return
	}

	// TODO: fetch paginated tags from the database using query.Page,
	// query.Limit, query.Sort and query.SortField.

	c.JSON(http.StatusOK, models.TagListResponse{
		Data: []models.Tag{},
		Meta: models.PaginationMeta{
			Total:      0,
			Page:       query.Page,
			Limit:      query.Limit,
			TotalPages: 0,
		},
	})
}

func (h *TagHandler) CreateTag(c *gin.Context) {
	var body models.TagBody
	if err := c.ShouldBindJSON(&body); err != nil {
		responses.ValidationError(c, err)
		return
	}

	// TODO: persist the new tag and obtain its generated ID.
	created := models.Tag{
		ID:      0,
		TagBody: body,
	}

	// TODO: set the Location header to the URL of the created resource,
	// e.g. c.Header("Location", fmt.Sprintf("/tasks/tags/%d", created.ID))

	c.JSON(http.StatusCreated, created)
}

func (h *TagHandler) GetTagByID(c *gin.Context) {
	var param models.TagIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		responses.ValidationError(c, err)
		return
	}

	// TODO: fetch the tag by param.TagID from the database.
	// If it does not exist, respond with httpresponse.NotFound(c, "...").

	c.JSON(http.StatusOK, models.Tag{ID: param.TagID})
}

func (h *TagHandler) UpdateTag(c *gin.Context) {
	var param models.TagIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		responses.ValidationError(c, err)
		return
	}

	var body models.TagBody
	if err := c.ShouldBindJSON(&body); err != nil {
		responses.ValidationError(c, err)
		return
	}

	// TODO: verify the tag with param.TagID exists
	// (httpresponse.NotFound(c, "...") if not), then persist the update.
	updated := models.Tag{
		ID:      param.TagID,
		TagBody: body,
	}

	c.JSON(http.StatusOK, updated)
}

func (h *TagHandler) DeleteTag(c *gin.Context) {
	var param models.TagIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		responses.ValidationError(c, err)
		return
	}

	// TODO: verify the tag with param.TagID exists
	// (httpresponse.NotFound(c, "...") if not), then delete it.

	c.Status(http.StatusNoContent)
}
