package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/chop1k/medods-test/internal/domain/models"
)

func tasksSortableFields() map[string]bool {
	return map[string]bool{
		"id":            true,
		"template_id":   true,
		"moved_task_id": true,
		"status":        true,
		"date":          true,
		"notes":         true,
		"starts_at":     true,
		"ends_at":       true,
	}
}

type TasksStorage struct {
	db *sql.DB
}

func NewTaskStorage(db *sql.DB) *TasksStorage {
	return &TasksStorage{
		db: db,
	}
}

func (s *TasksStorage) Transaction(fn func(tx *sql.Tx) (any, error)) (any, error) {
	tx, err := s.db.Begin()

	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	something, err := fn(tx)

	if err != nil {
		return nil, err
	}

	tx.Commit()

	return something, nil
}

func (s *TasksStorage) GetAll(tx *sql.Tx, page int, limit int, sortOrder string, sortField string) ([]models.Task, int, error) {
	fields := tasksSortableFields()

	_, ok := fields[sortField]

	if !ok {
		sortField = "id"
	}

	var query string

	if sortOrder == "asc" {
		query = fmt.Sprintf("select \"id\", \"template_id\", \"moved_task_id\", \"status\", \"notes\", \"date\", \"started_at\", \"ended_at\", \"deleted_at\", count(*) over() as total_count from \"app\".\"tasks\" order by %s asc limit $1 offset $2", sortField)
	} else {
		query = fmt.Sprintf("select \"id\", \"template_id\", \"moved_task_id\", \"status\", \"notes\", \"date\", \"started_at\", \"ended_at\", \"deleted_at\", count(*) over() as total_count from \"app\".\"tasks\" order by %s desc limit $1 offset $2", sortField)
	}

	var results *sql.Rows
	var err error

	if tx != nil {
		results, err = tx.Query(query, limit, (page-1)*limit)
	} else {
		results, err = s.db.Query(query, limit, (page-1)*limit)
	}

	if err != nil {
		return nil, 0, err
	}
	defer results.Close()

	tasks := []models.Task{}

	var count int

	for results.Next() {
		var task models.Task

		var date time.Time

		err = results.Scan(&task.ID, &task.TemplateID, &task.MovedId, &task.Status, &task.Notes, &date, &task.StartedAt, &task.EndedAt, &task.DeletedAt, &count)

		if err != nil {
			return nil, 0, err
		}

		formattedDate := date.Format("02-01-2006")

		task.Date = &formattedDate

		tasks = append(tasks, task)
	}

	return tasks, count, nil
}

func (s *TasksStorage) GetAllRunning(tx *sql.Tx, page int, limit int, sortOrder string, sortField string) ([]models.Task, int, error) {
	fields := tasksSortableFields()

	_, ok := fields[sortField]

	if !ok {
		sortField = "id"
	}

	var query string

	if sortOrder == "asc" {
		query = fmt.Sprintf("select \"id\", \"template_id\", \"moved_task_id\", \"status\", \"notes\", \"date\", \"started_at\", \"ended_at\", \"deleted_at\" count(*) over() as total_count from \"app\".\"tasks\" where \"status\" = 'running' order by %s asc limit $1 offset $2", sortField)
	} else {
		query = fmt.Sprintf("select \"id\", \"template_id\", \"moved_task_id\", \"status\", \"notes\", \"date\", \"started_at\", \"ended_at\", \"deleted_at\", count(*) over() as total_count from \"app\".\"tasks\" where \"status\" = 'running' order by %s desc limit $1 offset $2", sortField)
	}

	var results *sql.Rows
	var err error

	if tx != nil {
		results, err = tx.Query(query, limit, (page-1)*limit)
	} else {
		results, err = s.db.Query(query, limit, (page-1)*limit)
	}

	if err != nil {
		return nil, 0, err
	}
	defer results.Close()

	tasks := []models.Task{}

	var count int

	for results.Next() {
		var task models.Task

		var date time.Time

		err = results.Scan(&task.ID, &task.TemplateID, &task.MovedId, &task.Status, &task.Notes, &date, &task.StartedAt, &task.EndedAt, &task.DeletedAt, &count)

		if err != nil {
			return nil, 0, err
		}

		formattedDate := date.Format("02-01-2006")

		task.Date = &formattedDate

		tasks = append(tasks, task)
	}

	return tasks, count, nil
}

func (s *TasksStorage) GetById(tx *sql.Tx, id int) (*models.Task, error) {
	query := "select \"id\", \"template_id\", \"moved_task_id\", \"status\", \"notes\", \"date\", \"started_at\", \"ended_at\", \"deleted_at\" from \"app\".\"tasks\" where id = $1"

	var result *sql.Row

	if tx != nil {
		result = tx.QueryRow(query, id)
	} else {
		result = s.db.QueryRow(query, id)
	}

	var task models.Task

	var date time.Time

	err := result.Scan(&task.ID, &task.TemplateID, &task.MovedId, &task.Status, &task.Notes, &date, &task.StartedAt, &task.EndedAt, &task.DeletedAt)

	if err != nil {
		return nil, err
	}

	formattedDate := date.Format("02-01-2006")

	task.Date = &formattedDate

	return &task, nil
}

func (s *TasksStorage) Create(tx *sql.Tx, task models.TaskBody) (int, error) {
	date, err := time.Parse("02-01-2006", *task.Date)

	if err != nil {
		return 0, err
	}

	query := "insert into \"app\".\"tasks\" (\"template_id\", \"moved_task_id\", \"status\", \"notes\", \"date\", \"started_at\", \"ended_at\") values ($1, $2, $3, $4, $5, $6, $7) returning id"

	var result *sql.Row

	if tx != nil {
		result = tx.QueryRow(
			query,
			task.TemplateID,
			task.MovedId,
			task.Status,
			task.Notes,
			date.Format("2006-01-02"),
			task.StartedAt,
			task.EndedAt,
		)
	} else {
		result = s.db.QueryRow(
			query,
			task.TemplateID,
			task.MovedId,
			task.Status,
			task.Notes,
			date.Format("2006-01-02"),
			task.StartedAt,
			task.EndedAt,
		)
	}

	var id int

	err = result.Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *TasksStorage) UpdateById(tx *sql.Tx, id int, newTask models.TaskBody) (*models.Task, error) {
	date, err := time.Parse("02-01-2006", *newTask.Date)

	if err != nil {
		return nil, err
	}

	query := "update \"app\".\"tasks\" set \"template_id\" = $1, \"moved_task_id\" = $2, \"status\" = $3, \"notes\" = $4, \"date\" = $5, \"started_at\" = $6, \"ended_at\" = $7, \"deleted_at\" = $8 where id = $9 returning *"

	var result *sql.Row

	if tx != nil {
		result = tx.QueryRow(
			query,
			newTask.TemplateID,
			newTask.MovedId,
			newTask.Status,
			newTask.Notes,
			date,
			newTask.StartedAt,
			newTask.EndedAt,
			newTask.DeletedAt,
			id,
		)
	} else {
		result = s.db.QueryRow(
			query,
			newTask.TemplateID,
			newTask.MovedId,
			newTask.Status,
			newTask.Notes,
			date,
			newTask.StartedAt,
			newTask.EndedAt,
			newTask.DeletedAt,
			id,
		)
	}

	var task models.Task

	err = result.Scan(&task.ID, &task.TemplateID, &task.MovedId, &task.Status, &task.Notes, &date, &task.StartedAt, &task.EndedAt, &task.DeletedAt)

	if err != nil {
		return nil, err
	}

	formattedDate := date.Format("02-01-2006")

	task.Date = &formattedDate

	return &task, nil
}

func (s *TasksStorage) RemoveById(tx *sql.Tx, id int) error {
	query := "delete from \"app\".\"tasks\" where id = $1"

	var err error

	if tx != nil {
		_, err = tx.Exec(
			query,
			id,
		)
	} else {
		_, err = s.db.Exec(
			query,
			id,
		)
	}

	return err
}
