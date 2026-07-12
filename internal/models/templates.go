package models

import "time"

type TemplateSortableField string

const (
	TemplateSortByID          TemplateSortableField = "id"
	TemplateSortByName        TemplateSortableField = "name"
	TemplateSortByDescription TemplateSortableField = "description"
	TemplateSortByStartsAt    TemplateSortableField = "starts_at"
	TemplateSortByEndsAt      TemplateSortableField = "ends_at"
	TemplateSortByEnabled     TemplateSortableField = "enabled"
)

type SchedulingType string

const (
	SchedulingDaily    SchedulingType = "daily"
	SchedulingWeekly   SchedulingType = "weekly"
	SchedulingMonthly  SchedulingType = "monthly"
	SchedulingOneshot  SchedulingType = "oneshot"
	SchedulingEvenDays SchedulingType = "even-days"
	SchedulingOddDays  SchedulingType = "odd-days"
)

type Scheduling struct {
	Type SchedulingType `json:"type" binding:"required,oneof=daily weekly monthly oneshot even-days odd-days"`

	Include []string `json:"include,omitempty"`

	Exclude []string `json:"exclude,omitempty"`

	Dates []string `json:"dates,omitempty"`
}

type TemplateBody struct {
	Name        string    `json:"name" binding:"required,min=2,max=64"`
	Description string    `json:"description,omitempty" binding:"omitempty,min=2"`
	Tags        []int     `json:"tags,omitempty" binding:"omitempty,dive,min=1"`
	Enabled     bool      `json:"enabled" binding:"required"`
	StartsAt    time.Time `json:"starts_at" binding:"required"`
	EndsAt      time.Time `json:"ends_at" binding:"required"`

	Scheduling *Scheduling `json:"scheduling,omitempty"`
}

type Template struct {
	ID int `json:"id"`
	TemplateBody
}

type TemplateListResponse struct {
	Data []Template     `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type ListTemplatesQuery struct {
	Page      int    `form:"page,default=1" binding:"min=1"`
	Limit     int    `form:"limit,default=20" binding:"min=1,max=100"`
	Sort      string `form:"sort,default=asc" binding:"omitempty,oneof=asc desc"`
	SortField string `form:"sort-field,default=id" binding:"omitempty,oneof=id name description starts_at ends_at"`
}

type TemplateIDParam struct {
	TemplateID int `uri:"template_id" binding:"required,min=1"`
}
