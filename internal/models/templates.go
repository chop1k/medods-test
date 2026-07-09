package models

import "time"

// TemplateSortableField enumerates the fields templates can be sorted by,
// mirroring the `TemplateSortableField` schema.
type TemplateSortableField string

const (
	TemplateSortByID          TemplateSortableField = "id"
	TemplateSortByName        TemplateSortableField = "name"
	TemplateSortByDescription TemplateSortableField = "description"
	TemplateSortByStartsAt    TemplateSortableField = "starts_at"
	TemplateSortByEndsAt      TemplateSortableField = "ends_at"
)

// SchedulingType enumerates the discriminator values for the `scheduling`
// oneOf field on TemplateBody.
type SchedulingType string

const (
	SchedulingDaily    SchedulingType = "daily"
	SchedulingWeekly   SchedulingType = "weekly"
	SchedulingMonthly  SchedulingType = "monthly"
	SchedulingOneshot  SchedulingType = "oneshot"
	SchedulingEvenDays SchedulingType = "even-days"
	SchedulingOddDays  SchedulingType = "odd-days"
)

// Scheduling is a flattened representation of the oneOf
// TemplateDailyOptions / TemplateWeeklyOptions / TemplateMonthlyOptions /
// TemplateOneshotOptions / TemplateEvenDaysOptions / TemplateOddDaysOptions
// schemas, discriminated by Type.
//
// NOTE: which of Include/Exclude/Dates are required/applicable depends on
// Type (e.g. `include` is required for "weekly" and "monthly"). That
// cross-field validation is business logic and is left as a TODO in the
// handler layer rather than encoded here.
type Scheduling struct {
	Type SchedulingType `json:"type" binding:"required,oneof=daily weekly monthly oneshot even-days odd-days"`

	// Include is used by weekly (days of week) and monthly (days of month, 1-31) scheduling.
	Include []string `json:"include,omitempty"`

	// Exclude is used by daily/weekly/monthly/even-days/odd-days scheduling to
	// carve out specific days of week or wildcarded dates.
	Exclude []string `json:"exclude,omitempty"`

	// Dates is used by oneshot scheduling.
	Dates []string `json:"dates,omitempty"`
}

// TemplateBody mirrors the `TemplateBody` schema and is used as the request
// payload for create/update operations.
type TemplateBody struct {
	Name        string    `json:"name" binding:"required,min=2,max=64"`
	Description string    `json:"description" binding:"required,min=2"`
	Tags        []int64   `json:"tags,omitempty" binding:"omitempty,dive,min=1"`
	StartsAt    time.Time `json:"starts_at" binding:"required"`
	EndsAt      time.Time `json:"ends_at" binding:"required"`

	Scheduling *Scheduling `json:"scheduling,omitempty"`
}

// Template mirrors the `Template` schema (IdentifiableResource + TemplateBody).
type Template struct {
	ID int `json:"id"`
	TemplateBody
}

// TemplateListResponse is the envelope returned by GET /tasks/templates.
type TemplateListResponse struct {
	Data []Template     `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// ListTemplatesQuery captures the query parameters accepted by
// GET /tasks/templates.
type ListTemplatesQuery struct {
	Page      int    `form:"page,default=1" binding:"min=1"`
	Limit     int    `form:"limit,default=20" binding:"min=1,max=100"`
	Sort      string `form:"sort,default=asc" binding:"omitempty,oneof=asc desc"`
	SortField string `form:"sort-field,default=id" binding:"omitempty,oneof=id name description starts_at ends_at"`
}

// TemplateIDParam captures the `template_id` path parameter shared by
// GET/PUT/DELETE /tasks/templates/{template_id}.
type TemplateIDParam struct {
	TemplateID int `uri:"template_id" binding:"required,min=1"`
}
