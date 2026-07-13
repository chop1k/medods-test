package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/chop1k/medods-test/internal/domain/models"
)

// Error writes a generic ErrorResponse with the given status/title/detail.
func Error(c *gin.Context, status int, title, detail string) {
	c.JSON(status, models.ErrorResponse{
		Title:  title,
		Status: status,
		Detail: detail,
	})
}

// ValidationError writes a 400 ValidationErrorResponse built from a Gin
// binding error. Translating the raw binding error into structured
// per-field messages is left as a TODO - for now the raw error text is
// surfaced in Detail so the boilerplate is still useful during development.
func ValidationError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, models.ValidationErrorResponse{
		ErrorResponse: models.ErrorResponse{
			Title:  "Bad Request",
			Status: http.StatusBadRequest,
			Detail: err.Error(),
		},
		// TODO: populate Errors with per-field {field, message} entries once
		// validation error translation is implemented.
		Errors: nil,
	})
}

// NotFound writes a 404 ErrorResponse.
func NotFound(c *gin.Context, detail string) {
	Error(c, http.StatusNotFound, "Not Found", detail)
}

// InternalError writes a 500 ErrorResponse.
func InternalError(c *gin.Context, detail string) {
	Error(c, http.StatusInternalServerError, "Internal Server Error", detail)
}
