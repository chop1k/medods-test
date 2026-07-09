package models

import "time"

// TaskSortableField enumerates the fields tasks can be sorted by, mirroring
// the `TaskSortableField` schema.
type TaskSortableField string

const (
	TaskSortByID         TaskSortableField = "id"
	TaskSortByTemplateID TaskSortableField = "template_id"
	TaskSortByStatus     TaskSortableField = "status"
	TaskSortByNotes      TaskSortableField = "notes"
	TaskSortByStartedAt  TaskSortableField = "started_at"
	TaskSortByEndedAt    TaskSortableField = "ended_at"
)

// TaskStatus enumerates the `status` field values on TaskBody.
type TaskStatus string

const (
	TaskStatusPending  TaskStatus = "pending"
	TaskStatusRunning  TaskStatus = "running"
	TaskStatusFinished TaskStatus = "finished"
)

// TaskBody mirrors the `TaskBody` schema and is used as the request payload
// for create/update operations.
type TaskBody struct {
	TemplateID int64      `json:"template_id" binding:"required,min=1"`
	Status     TaskStatus `json:"status,omitempty" binding:"omitempty,oneof=pending running finished"`
	Notes      string     `json:"notes,omitempty" binding:"omitempty,min=2"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	EndedAt    *time.Time `json:"ended_at,omitempty"`
}

// Task mirrors the `Task` schema (IdentifiableResource + TaskBody).
type Task struct {
	ID int64 `json:"id"`
	TaskBody
}

// TaskListResponse is the envelope returned by GET /tasks.
type TaskListResponse struct {
	Data []Task         `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// ListTasksQuery captures the query parameters accepted by GET /tasks.
type ListTasksQuery struct {
	Page      int    `form:"page,default=1" binding:"min=1"`
	Limit     int    `form:"limit,default=20" binding:"min=1,max=100"`
	Sort      string `form:"sort,default=asc" binding:"omitempty,oneof=asc desc"`
	SortField string `form:"sort-field,default=id" binding:"omitempty,oneof=id template_id status notes started_at ended_at"`
}

// TaskIDParam captures the `task_id` path parameter shared by
// GET/PUT/DELETE /tasks/{task_id}.
type TaskIDParam struct {
	TaskID int64 `uri:"task_id" binding:"required,min=1"`
}
