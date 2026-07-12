package repository

import (
	"database/sql"

	"github.com/chop1k/medods-test/internal/models"
)

// taskColumns lists the columns that map to models.TaskBody. The
// "app"."tasks" table also has "moved_task_id" and "date" columns that
// aren't represented on the model yet, so they're intentionally left out
// here rather than guessed at - add them (and the matching model fields)
// once that mapping is decided.
const taskColumns = `"id", "template_id", "status", "notes", "started_at", "ended_at"`

type TasksStorage struct {
	db *sql.DB
}

func NewTaskStorage(db *sql.DB) *TasksStorage {
	return &TasksStorage{
		db: db,
	}
}

func (s *TasksStorage) GetAll(page int, limit int) ([]models.Task, error) {
	results, err := s.db.Query(
		"select "+taskColumns+" from \"app\".\"tasks\" limit $1 offset $2",
		limit, (page-1)*limit,
	)

	if err != nil {
		return nil, err
	}
	defer results.Close()

	tasks := []models.Task{}

	for results.Next() {
		var task models.Task

		err = results.Scan(&task.ID, &task.TemplateID, &task.Status, &task.Notes, &task.StartedAt, &task.EndedAt)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *TasksStorage) GetById(id int) (*models.Task, error) {
	result := s.db.QueryRow(
		"select "+taskColumns+" from \"app\".\"tasks\" where id = $1",
		id,
	)

	var task models.Task

	err := result.Scan(&task.ID, &task.TemplateID, &task.Status, &task.Notes, &task.StartedAt, &task.EndedAt)

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *TasksStorage) Create(task models.TaskBody) (int, error) {
	result := s.db.QueryRow(
		"insert into \"app\".\"tasks\" (\"template_id\", \"status\", \"notes\", \"started_at\", \"ended_at\") values ($1, $2, $3, $4, $5) returning id",
		task.TemplateID,
		task.Status,
		task.Notes,
		task.StartedAt,
		task.EndedAt,
	)

	var id int

	err := result.Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *TasksStorage) UpdateById(id int, newTask models.TaskBody) (*models.Task, error) {
	result := s.db.QueryRow(
		"update \"app\".\"tasks\" set \"template_id\" = $1, \"status\" = $2, \"notes\" = $3, \"started_at\" = $4, \"ended_at\" = $5 where id = $6 returning "+taskColumns,
		newTask.TemplateID,
		newTask.Status,
		newTask.Notes,
		newTask.StartedAt,
		newTask.EndedAt,
		id,
	)

	var task models.Task

	err := result.Scan(&task.ID, &task.TemplateID, &task.Status, &task.Notes, &task.StartedAt, &task.EndedAt)

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *TasksStorage) RemoveById(id int) error {
	_, err := s.db.Exec(
		"delete from \"app\".\"tasks\" where id = $1",
		id,
	)

	return err
}