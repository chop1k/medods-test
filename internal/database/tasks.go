package database

import "github.com/chop1k/medods-test/internal/models"

type TasksStorage struct {
}

func (s *TasksStorage) GetAll(page int, limit int) {
}

func (s *TasksStorage) GetById(id int) {
}

func (s *TasksStorage) Create(task models.TaskBody) {
}

func (s *TasksStorage) CreateBulk(tasks []models.TaskBody) ([]int, error) {
	panic("not implemented")
}

func (s *TasksStorage) UpdateById(id int, task models.TaskBody) {
}

func (s *TasksStorage) RemoveById(id int) {
}
