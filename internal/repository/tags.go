package repository

import (
	"database/sql"
	"fmt"

	"github.com/chop1k/medods-test/internal/domain/models"
)

// tagColumns lists the columns that map to models.TagBody. The
// "app"."tags" table also has a "description" column that isn't
// represented on the model yet, so it's intentionally left out here rather
// than guessed at - add it (and the matching model field) once that
// mapping is decided.
const tagColumns = ``

func tagsSortableFields() map[string]bool {
	return map[string]bool{
		"id":          true,
		"name":        true,
		"description": true,
	}
}

type TagsStorage struct {
	db *sql.DB
}

func NewTagStorage(db *sql.DB) *TagsStorage {
	return &TagsStorage{
		db: db,
	}
}

func (s *TagsStorage) GetAll(page int, limit int, sortOrder string, sortField string) ([]models.Tag, int, error) {
	fields := tagsSortableFields()

	_, ok := fields[sortField]

	if !ok {
		sortField = "id"
	}

	var query string

	if sortOrder == "asc" {
		query = fmt.Sprintf("select \"id\", \"name\", \"description\", \"type\", count(*) over() as total_count from \"app\".\"tags\" order by %s asc limit $1 offset $2", sortField)
	} else {
		query = fmt.Sprintf("select \"id\", \"name\", \"description\", \"type\", count(*) over() as total_count from \"app\".\"tags\" order by %s desc limit $1 offset $2", sortField)
	}

	results, err := s.db.Query(query, limit, (page-1)*limit)

	if err != nil {
		return nil, 0, err
	}
	defer results.Close()

	tags := []models.Tag{}

	var count int

	for results.Next() {
		var tag models.Tag

		err = results.Scan(&tag.ID, &tag.Name, &tag.Description, &tag.Type, &count)

		if err != nil {
			return nil, 0, err
		}

		tags = append(tags, tag)
	}

	return tags, count, nil
}

func (s *TagsStorage) GetById(id int) (*models.Tag, error) {
	result := s.db.QueryRow(
		"select \"id\", \"name\", \"description\", \"type\" from \"app\".\"tags\" where id = $1",
		id,
	)

	var tag models.Tag

	err := result.Scan(&tag.ID, &tag.Name, &tag.Description, &tag.Type)

	if err != nil {
		return nil, err
	}

	return &tag, nil
}

func (s *TagsStorage) Create(tag models.TagBody) (int, error) {
	result := s.db.QueryRow(
		"insert into \"app\".\"tags\" (\"name\", \"description\", \"type\") values ($1, $2, 'user-defined') returning id",
		tag.Name,
		tag.Description,
	)

	var id int

	err := result.Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *TagsStorage) RemoveById(id int) error {
	_, err := s.db.Exec(
		"delete from \"app\".\"tags\" where id = $1",
		id,
	)

	return err
}
