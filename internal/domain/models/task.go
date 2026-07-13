package models

type TaskSortableField string

const (
	TaskSortByID          TaskSortableField = "id"
	TaskSortByTemplateID  TaskSortableField = "template_id"
	TaskSortByMovedTaskID TaskSortableField = "moved_task_id"
	TaskSortByDate        TaskSortableField = "date"
	TaskSortByStatus      TaskSortableField = "status"
	TaskSortByNotes       TaskSortableField = "notes"
	TaskSortByStartedAt   TaskSortableField = "started_at"
	TaskSortByEndedAt     TaskSortableField = "ended_at"
)

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusFinished  TaskStatus = "finished"
	TaskStatusCancelled TaskStatus = "cancelled"
	TaskStatusMoved     TaskStatus = "moved"
	TaskStatusOverdue   TaskStatus = "overdue"
)

type TaskBody struct {
	TemplateID *int       `json:"template_id" binding:"required,min=1"`
	MovedId    *int       `json:"moved_to" binding:"omitempty"`
	Status     TaskStatus `json:"status,omitempty" binding:"omitempty,oneof=pending running finished cancelled moved overdue"`
	Notes      *string    `json:"notes,omitempty" binding:"omitempty,min=2"`
	Date       *string    `json:"date" binding:"required"`
	StartedAt  *string    `json:"started_at,omitempty"`
	EndedAt    *string    `json:"ended_at,omitempty"`
	DeletedAt  *string    `json:"deleted_at"`
}

type RunningTaskBody struct {
	Status    *TaskStatus `json:"status,omitempty" binding:"required,oneof=running"`
	StartedAt *string     `json:"started_at" binding:"required"`
}

type FinishedTaskBody struct {
	Status  *TaskStatus `json:"status" binding:"required,oneof=finished"`
	EndedAt *string     `json:"ended_at" binding:"required"`
}

type CancelledTaskBody struct {
	Status    *TaskStatus `json:"status" binding:"required,oneof=cancelled"`
	StartedAt *string     `json:"started_at,omitempty"`
	EndedAt   *string     `json:"ended_at,omitempty"`
}

type MovedTaskBody struct {
	Status *TaskStatus `json:"status" binding:"required,oneof=moved"`
	Date   *string     `json:"date" binding:"required"`
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
	SortField string `form:"sort-field,default=id" binding:"omitempty,oneof=id template_id moved_task_id status notes date started_at ended_at"`
}

type TaskIDParam struct {
	TaskID int `uri:"task_id" binding:"required,min=1"`
}
