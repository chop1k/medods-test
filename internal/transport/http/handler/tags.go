package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/chop1k/medods-test/internal/models"
	"github.com/chop1k/medods-test/internal/repository"
)

type TagHandler struct {
	repository *repository.TagsStorage
}

func NewTagHandler(repository *repository.TagsStorage) *TagHandler {
	return &TagHandler{
		repository: repository,
	}
}

func (h *TagHandler) GetTags(c *gin.Context) {
	var query models.ListTagsQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		ValidationError(c, err)

		return
	}

	tags, err := h.repository.GetAll(query.Page, query.Limit)

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, models.TagListResponse{
		Data: tags,
		Meta: models.PaginationMeta{
			Total:      len(tags),
			Page:       query.Page,
			Limit:      query.Limit,
			TotalPages: 0,
		},
	})
}

func (h *TagHandler) CreateTag(c *gin.Context) {
	var body models.TagBody

	if err := c.ShouldBindJSON(&body); err != nil {
		ValidationError(c, err)

		return
	}

	id, err := h.repository.Create(body)

	if err != nil {
		panic(err)
	}

	created := models.Tag{
		ID:      id,
		TagBody: body,
	}

	c.Header("Location", "/v1/grouping/tags/"+strconv.Itoa(id))
	c.JSON(http.StatusCreated, created)
}

func (h *TagHandler) GetTagByID(c *gin.Context) {
	var param models.TagIDParam

	if err := c.ShouldBindUri(&param); err != nil {
		ValidationError(c, err)

		return
	}

	tag, err := h.repository.GetById(param.TagID)

	if err != nil {
		NotFound(c, err.Error())

		return
	}

	c.JSON(http.StatusOK, tag)
}

func (h *TagHandler) UpdateTag(c *gin.Context) {
	var param models.TagIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		ValidationError(c, err)
		return
	}

	var body models.TagBody
	if err := c.ShouldBindJSON(&body); err != nil {
		ValidationError(c, err)
		return
	}

	updated, err := h.repository.UpdateById(param.TagID, body)

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, updated)
}

func (h *TagHandler) DeleteTag(c *gin.Context) {
	var param models.TagIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		ValidationError(c, err)
		return
	}

	err := h.repository.RemoveById(param.TagID)

	if err != nil {
		panic(err)
	}

	c.Status(http.StatusNoContent)
}
