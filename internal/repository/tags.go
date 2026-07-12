package repository

import (
	"database/sql"

	"github.com/chop1k/medods-test/internal/models"
)

// tagColumns lists the columns that map to models.TagBody. The
// "app"."tags" table also has a "description" column that isn't
// represented on the model yet, so it's intentionally left out here rather
// than guessed at - add it (and the matching model field) once that
// mapping is decided.
const tagColumns = `"id", "name", "type", "deleted_at"`

type TagsStorage struct {
	db *sql.DB
}

func NewTagStorage(db *sql.DB) *TagsStorage {
	return &TagsStorage{
		db: db,
	}
}

func (s *TagsStorage) GetAll(page int, limit int) ([]models.Tag, error) {
	results, err := s.db.Query(
		"select "+tagColumns+" from \"app\".\"tags\" limit $1 offset $2",
		limit, (page-1)*limit,
	)

	if err != nil {
		return nil, err
	}
	defer results.Close()

	tags := []models.Tag{}

	for results.Next() {
		var tag models.Tag

		err = results.Scan(&tag.ID, &tag.Name, &tag.Type, &tag.DeletedAt)

		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (s *TagsStorage) GetById(id int) (*models.Tag, error) {
	result := s.db.QueryRow(
		"select "+tagColumns+" from \"app\".\"tags\" where id = $1",
		id,
	)

	var tag models.Tag

	err := result.Scan(&tag.ID, &tag.Name, &tag.Type, &tag.DeletedAt)

	if err != nil {
		return nil, err
	}

	return &tag, nil
}

func (s *TagsStorage) Create(tag models.TagBody) (int, error) {
	result := s.db.QueryRow(
		"insert into \"app\".\"tags\" (\"name\", \"type\", \"deleted_at\") values ($1, $2, $3) returning id",
		tag.Name,
		tag.Type,
		tag.DeletedAt,
	)

	var id int

	err := result.Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *TagsStorage) UpdateById(id int, newTag models.TagBody) (*models.Tag, error) {
	result := s.db.QueryRow(
		"update \"app\".\"tags\" set \"name\" = $1, \"type\" = $2, \"deleted_at\" = $3 where id = $4 returning "+tagColumns,
		newTag.Name,
		newTag.Type,
		newTag.DeletedAt,
		id,
	)

	var tag models.Tag

	err := result.Scan(&tag.ID, &tag.Name, &tag.Type, &tag.DeletedAt)

	if err != nil {
		return nil, err
	}

	return &tag, nil
}

func (s *TagsStorage) RemoveById(id int) error {
	_, err := s.db.Exec(
		"delete from \"app\".\"tags\" where id = $1",
		id,
	)

	return err
}
