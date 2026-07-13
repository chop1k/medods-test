package models

type ErrorResponse struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title,omitempty"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

type ValidationErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrorResponse struct {
	ErrorResponse
	Errors []ValidationErrorDetail `json:"errors,omitempty"`
}

type PaginationMeta struct {
	Total      int `json:"total"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
}

// ListQuery captures the pagination/sorting query parameters shared by all
// collection endpoints (`page`, `limit`, `sort`, `sort-field`).
//
// SortField is intentionally typed as `string` here; handlers that need a
// restricted enum (e.g. TemplateSortableField) should validate it via
// `binding:"oneof=..."` on an embedding struct or a manual check.
type ListQuery struct {
	Page      int    `form:"page,default=1" binding:"min=1"`
	Limit     int    `form:"limit,default=20" binding:"min=1,max=100"`
	Sort      string `form:"sort,default=asc" binding:"omitempty,oneof=asc desc"`
	SortField string `form:"sort-field"`
}
