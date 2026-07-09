package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/chop1k/medods-test/internal/database"
	"github.com/chop1k/medods-test/internal/models"
	"github.com/chop1k/medods-test/internal/responses"
)

type TemplateHandler struct {
	repository *database.TemplatesStorage
}

func NewTemplateHandler(storage *database.TemplatesStorage) *TemplateHandler {
	return &TemplateHandler{
		repository: storage,
	}
}

func (h *TemplateHandler) GetTemplates(c *gin.Context) {
	var query models.ListTemplatesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		responses.ValidationError(c, err)
		return
	}

	// TODO: fetch paginated templates from the database using query.Page,
	// query.Limit, query.Sort and query.SortField.

	c.JSON(http.StatusOK, models.TemplateListResponse{
		Data: []models.Template{},
		Meta: models.PaginationMeta{
			Total:      0,
			Page:       query.Page,
			Limit:      query.Limit,
			TotalPages: 0,
		},
	})
}

func (h *TemplateHandler) CreateTemplate(c *gin.Context) {
	var body models.TemplateBody
	if err := c.ShouldBindJSON(&body); err != nil {
		responses.ValidationError(c, err)
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

	// TODO: set the Location header to the URL of the created resource,
	// e.g. c.Header("Location", fmt.Sprintf("/tasks/templates/%d", created.ID))

	c.Header("Location", "/v1/tasks/task"+strconv.Itoa(id))
	c.JSON(http.StatusCreated, created)
}

func (h *TemplateHandler) GetTemplateByID(c *gin.Context) {
	var param models.TemplateIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		responses.ValidationError(c, err)
		return
	}

	// TODO: fetch the template by param.TemplateID from the database.
	// If it does not exist, respond with httpresponse.NotFound(c, "...").

	c.JSON(http.StatusOK, models.Template{ID: param.TemplateID})
}

func (h *TemplateHandler) UpdateTemplate(c *gin.Context) {
	var param models.TemplateIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		responses.ValidationError(c, err)
		return
	}

	var body models.TemplateBody
	if err := c.ShouldBindJSON(&body); err != nil {
		responses.ValidationError(c, err)
		return
	}

	// TODO: verify the template with param.TemplateID exists
	// (httpresponse.NotFound(c, "...") if not), then persist the update.
	updated := models.Template{
		ID:           param.TemplateID,
		TemplateBody: body,
	}

	c.JSON(http.StatusOK, updated)
}

// DeleteTemplate handles DELETE /tasks/templates/{template_id}
// (operationId: deleteTemplate).
func (h *TemplateHandler) DeleteTemplate(c *gin.Context) {
	var param models.TemplateIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		responses.ValidationError(c, err)
		return
	}

	// TODO: verify the template with param.TemplateID exists
	// (httpresponse.NotFound(c, "...") if not), then delete it.

	c.Status(http.StatusNoContent)
}
