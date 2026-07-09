package database

import "github.com/chop1k/medods-test/internal/models"

type TasksStorage struct {
}

func (s *TasksStorage) GetAll(page int, limit int) {
}

func (s *TasksStorage) GetById(id int) {
}

func (s *TasksStorage) Create(template models.TaskBody) {
}

func (s *TasksStorage) UpdateById(id int, template models.TaskBody) {
}

func (s *TasksStorage) RemoveById(id int) {
}
