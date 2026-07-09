package models

import "time"

type TaskSortableField string

const (
	TaskSortByID         TaskSortableField = "id"
	TaskSortByTemplateID TaskSortableField = "template_id"
	TaskSortByDate       TaskSortableField = "date"
	TaskSortByStatus     TaskSortableField = "status"
	TaskSortByNotes      TaskSortableField = "notes"
	TaskSortByStartedAt  TaskSortableField = "started_at"
	TaskSortByEndedAt    TaskSortableField = "ended_at"
)

type TaskStatus string

const (
	TaskStatusPending  TaskStatus = "pending"
	TaskStatusRunning  TaskStatus = "running"
	TaskStatusFinished TaskStatus = "finished"
)

type TaskBody struct {
	TemplateID int        `json:"template_id" binding:"required,min=1"`
	Status     TaskStatus `json:"status,omitempty" binding:"omitempty,oneof=pending running finished"`
	Notes      *string    `json:"notes,omitempty" binding:"omitempty,min=2"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	EndedAt    *time.Time `json:"ended_at,omitempty"`
}

type Task struct {
	ID int `json:"id"`
	TaskBody
}

type TaskListResponse struct {
	Data []Task         `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type ListTasksQuery struct {
	Page      int    `form:"page,default=1" binding:"min=1"`
	Limit     int    `form:"limit,default=20" binding:"min=1,max=100"`
	Sort      string `form:"sort,default=asc" binding:"omitempty,oneof=asc desc"`
	SortField string `form:"sort-field,default=id" binding:"omitempty,oneof=id template_id status notes started_at ended_at"`
}

type TaskIDParam struct {
	TaskID int `uri:"task_id" binding:"required,min=1"`
}
