package handler

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/chop1k/medods-test/internal/domain/models"
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

	fmt.Printf("%s %s\n", query.Sort, query.SortField)

	templates, count, err := h.repository.GetAll(query.Page, query.Limit, query.Sort, query.SortField)

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, models.TemplateListResponse{
		Data: templates,
		Meta: models.PaginationMeta{
			Total:      len(templates),
			Page:       query.Page,
			Limit:      query.Limit,
			TotalPages: int(math.Ceil(float64(count) / float64(query.Limit))),
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

	template, err := h.repository.GetById(param.TemplateID)

	if err != nil {
		panic(err)
	}

	if template.DeletedAt != nil {
		NotFound(c, "not found")

		return
	}

	var body models.TemplateUpdateBody
	if err := c.ShouldBindJSON(&body); err != nil {
		ValidationError(c, err)

		return
	}

	template, err = h.repository.UpdateById(param.TemplateID, body)

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

	template, err := h.repository.GetById(param.TemplateID)

	if err != nil {
		panic(err)
	}

	if template.DeletedAt != nil {
		NotFound(c, "not found")

		return
	}

	enabled := false
	deletedAt := time.Now().Format(time.RFC3339)

	updated := models.TemplateUpdateBody{
		Enabled:   &enabled,
		DeletedAt: &deletedAt,
	}

	_, err = h.repository.UpdateById(param.TemplateID, updated)

	if err != nil {
		panic(err)
	}

	c.Status(http.StatusNoContent)
}
