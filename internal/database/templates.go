package database

import (
	"database/sql"
	"encoding/json"

	"github.com/chop1k/medods-test/internal/models"
)

type TemplatesStorage struct {
	db *sql.DB
}

func NewTemplateStorage(db *sql.DB) *TemplatesStorage {
	return &TemplatesStorage{
		db: db,
	}
}

func (s *TemplatesStorage) GetAll(page int, limit int) ([]models.Template, error) {
	results, err := s.db.Query("select * from \"app\".\"templates\" limit $1 offset $2", limit, (page-1)*limit)

	if err != nil {
		return nil, err
	}
	defer results.Close()

	templates := []models.Template{}

	for results.Next() {
		var template models.Template
		var schedulingRaw []byte

		err = results.Scan(&template.ID, &template.Name, &template.Description, &template.StartsAt, &template.EndsAt, &template.Enabled, &schedulingRaw)

		if err != nil {
			return nil, err
		}

		var scheduling models.Scheduling

		err = json.Unmarshal(schedulingRaw, &scheduling)

		if err != nil {
			return nil, err
		}

		template.Scheduling = &scheduling

		templates = append(templates, template)
	}

	return templates, nil
}

func (s *TemplatesStorage) GetById(id int) (*models.Template, error) {
	results := s.db.QueryRow("select * from \"app\".\"templates\" where id = $1", id)

	var template models.Template
	var schedulingRaw []byte

	err := results.Scan(&template.ID, &template.Name, &template.Description, &template.StartsAt, &template.EndsAt, &template.Enabled, &schedulingRaw)

	if err != nil {
		return nil, err
	}

	var scheduling models.Scheduling

	err = json.Unmarshal(schedulingRaw, &scheduling)

	if err != nil {
		return nil, err
	}

	template.Scheduling = &scheduling

	return &template, nil
}

func (s *TemplatesStorage) Create(template models.TemplateBody) (int, error) {
	scheduling, err := json.Marshal(template.Scheduling)

	if err != nil {
		return 0, err
	}

	results := s.db.QueryRow(
		"insert into \"app\".\"templates\" (\"name\", \"description\", \"starts_at\", \"ends_at\", \"enabled\", \"scheduling\") values ($1, $2, $3, $4, $5, $6) returning id",
		template.Name,
		template.Description,
		template.StartsAt,
		template.EndsAt,
		template.Enabled,
		scheduling,
	)

	var id int

	err = results.Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *TemplatesStorage) UpdateById(id int, newTemplate models.TemplateBody) (*models.Template, error) {
	newScheduling, err := json.Marshal(newTemplate.Scheduling)

	if err != nil {
		return nil, err
	}

	result := s.db.QueryRow(
		"update \"app\".\"templates\" set \"name\" = $1, \"description\" = $2, \"starts_at\" = $3, \"ends_at\" = $4, \"enabled\" = $5, \"scheduling\" = $6 where id = $7 returning *",
		newTemplate.Name,
		newTemplate.Description,
		newTemplate.StartsAt,
		newTemplate.EndsAt,
		newTemplate.Enabled,
		newScheduling,
		id,
	)

	var template models.Template
	var schedulingRaw []byte

	err = result.Scan(&template.ID, &template.Name, &template.Description, &template.StartsAt, &template.EndsAt, &schedulingRaw)

	if err != nil {
		return nil, err
	}

	var scheduling models.Scheduling

	err = json.Unmarshal(schedulingRaw, &scheduling)

	if err != nil {
		return nil, err
	}

	template.Scheduling = &scheduling

	return &template, nil
}

func (s *TemplatesStorage) RemoveById(id int) error {
	_, err := s.db.Exec(
		"delete from \"app\".\"templates\" where id = $1",
		id,
	)

	return err
}

func (s *TemplatesStorage) GetAllDaily(page int, limit int) ([]models.Template, error) {
	results, err := s.db.Query("select * from \"app\".\"templates\" where scheduling ->> 'type' = 'daily' limit $1 offset $2", limit, (page-1)*limit)

	if err != nil {
		return nil, err
	}
	defer results.Close()

	templates := []models.Template{}

	for results.Next() {
		var template models.Template
		var schedulingRaw []byte

		err = results.Scan(&template.ID, &template.Name, &template.Description, &template.StartsAt, &template.EndsAt, &template.Enabled, &schedulingRaw)

		if err != nil {
			return nil, err
		}

		var scheduling models.Scheduling

		err = json.Unmarshal(schedulingRaw, &scheduling)

		if err != nil {
			return nil, err
		}

		template.Scheduling = &scheduling

		templates = append(templates, template)
	}

	return templates, nil
}
