package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/chop1k/medods-test/internal/models"
	"github.com/chop1k/medods-test/internal/repository"
)

type TemplateHandler struct {
	repository *repository.TemplatesStorage
}

func NewTemplateHandler(storage *repository.TemplatesStorage) *TemplateHandler {
	return &TemplateHandler{
		repository: storage,
	}
}

func (h *TemplateHandler) GetTemplates(c *gin.Context) {
	var query models.ListTemplatesQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		ValidationError(c, err)

		return
	}

	templates, err := h.repository.GetAll(query.Page, query.Limit)

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, models.TemplateListResponse{
		Data: templates,
		Meta: models.PaginationMeta{
			Total:      len(templates),
			Page:       query.Page,
			Limit:      query.Limit,
			TotalPages: 0,
		},
	})
}

func (h *TemplateHandler) CreateTemplate(c *gin.Context) {
	var body models.TemplateBody

	if err := c.ShouldBindJSON(&body); err != nil {
		ValidationError(c, err)
		return
	}

	id, err := h.repository.Create(body)

	if err != nil {
		panic(err)
	}

	created := models.Template{
		ID:           id,
		TemplateBody: body,
	}

	c.Header("Location", "/v1/tasks/template/"+strconv.Itoa(id))
	c.JSON(http.StatusCreated, created)
}

func (h *TemplateHandler) GetTemplateByID(c *gin.Context) {
	var param models.TemplateIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		ValidationError(c, err)

		return
	}

	template, err := h.repository.GetById(param.TemplateID)

	if err != nil {
		NotFound(c, err.Error())

		return
	}

	c.JSON(http.StatusOK, template)
}

func (h *TemplateHandler) UpdateTemplate(c *gin.Context) {
	var param models.TemplateIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		ValidationError(c, err)

		return
	}

	var body models.TemplateUpdateBody
	if err := c.ShouldBindJSON(&body); err != nil {
		ValidationError(c, err)

		return
	}

	template, err := h.repository.UpdateById(param.TemplateID, body)

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, template)
}

func (h *TemplateHandler) DeleteTemplate(c *gin.Context) {
	var param models.TemplateIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		ValidationError(c, err)

		return
	}

	err := h.repository.RemoveById(param.TemplateID)

	if err != nil {
		panic(err)
	}

	c.Status(http.StatusNoContent)
}
