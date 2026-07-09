package database

import "github.com/chop1k/medods-test/internal/models"

type TagsStorage struct {
}

func (s *TagsStorage) GetAll(page int, limit int) {
}

func (s *TagsStorage) GetById(id int) {
}

func (s *TagsStorage) Create(template models.TagBody) {
}

func (s *TagsStorage) UpdateById(id int, template models.TagBody) {
}

func (s *TagsStorage) RemoveById(id int) {
}
