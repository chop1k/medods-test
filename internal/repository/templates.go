package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/chop1k/medods-test/internal/domain/models"
)

func templatesSortableFields() map[string]bool {
	return map[string]bool{
		"id":          true,
		"name":        true,
		"description": true,
		"starts_at":   true,
		"ends_at":     true,
	}
}

type TemplatesStorage struct {
	db *sql.DB
}

func NewTemplateStorage(db *sql.DB) *TemplatesStorage {
	return &TemplatesStorage{
		db: db,
	}
}

func (s *TemplatesStorage) GetAll(page int, limit int, sortOrder string, sortField string) ([]models.Template, int, error) {
	fields := templatesSortableFields()

	_, ok := fields[sortField]

	if !ok {
		sortField = "id"
	}

	var query string

	if sortOrder == "asc" {
		query = fmt.Sprintf("select *, count(*) over() as total_count from \"app\".\"templates\" where \"deleted_at\" is null order by %s asc limit $1 offset $2", sortField)
	} else {
		query = fmt.Sprintf("select *, count(*) over() as total_count from \"app\".\"templates\" where \"deleted_at\" is null order by %s desc limit $1 offset $2", sortField)
	}

	results, err := s.db.Query(query, limit, (page-1)*limit)

	if err != nil {
		return nil, 0, err
	}
	defer results.Close()

	templates := []models.Template{}

	var count int

	for results.Next() {
		var template models.Template
		var schedulingRaw []byte

		err = results.Scan(&template.ID, &template.Name, &template.Description, &template.StartsAt, &template.EndsAt, &template.Enabled, &schedulingRaw, &template.DeletedAt, &count)

		if err != nil {
			return nil, 0, err
		}

		var scheduling models.Scheduling

		err = json.Unmarshal(schedulingRaw, &scheduling)

		if err != nil {
			return nil, 0, err
		}

		template.Scheduling = &scheduling

		templates = append(templates, template)
	}

	return templates, count, nil
}

func (s *TemplatesStorage) GetById(id int) (*models.Template, error) {
	results := s.db.QueryRow("select * from \"app\".\"templates\" where id = $1 and \"deleted_at\" is null", id)

	var template models.Template
	var schedulingRaw []byte

	err := results.Scan(&template.ID, &template.Name, &template.Description, &template.StartsAt, &template.EndsAt, &template.Enabled, &schedulingRaw, &template.DeletedAt)

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

func (s *TemplatesStorage) UpdateById(id int, newTemplate models.TemplateUpdateBody) (*models.Template, error) {
	result := s.db.QueryRow(
		"update \"app\".\"templates\" set \"enabled\" = $1, \"deleted_at\" = $2 where id = $3 returning *",
		newTemplate.Enabled,
		newTemplate.DeletedAt,
		id,
	)

	var template models.Template
	var schedulingRaw []byte

	err := result.Scan(&template.ID, &template.Name, &template.Description, &template.StartsAt, &template.EndsAt, &template.Enabled, &schedulingRaw, &template.DeletedAt)

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

func (s *TemplatesStorage) GetAllDaily(page int, limit int) ([]models.Template, error) {
	results, err := s.db.Query("select * from \"app\".\"templates\" where scheduling ->> 'type' = 'daily' and \"deleted_at\" is null limit $1 offset $2", limit, (page-1)*limit)

	if err != nil {
		return nil, err
	}
	defer results.Close()

	templates := []models.Template{}

	for results.Next() {
		var template models.Template
		var schedulingRaw []byte

		err = results.Scan(&template.ID, &template.Name, &template.Description, &template.StartsAt, &template.EndsAt, &template.Enabled, &schedulingRaw, &template.DeletedAt)

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

func (s *TemplatesStorage) GetAllWeekly(page int, limit int) ([]models.Template, error) {
	results, err := s.db.Query("select * from \"app\".\"templates\" where scheduling ->> 'type' = 'weekly' and \"deleted_at\" is null limit $1 offset $2", limit, (page-1)*limit)

	if err != nil {
		return nil, err
	}
	defer results.Close()

	templates := []models.Template{}

	for results.Next() {
		var template models.Template
		var schedulingRaw []byte

		err = results.Scan(&template.ID, &template.Name, &template.Description, &template.StartsAt, &template.EndsAt, &template.Enabled, &schedulingRaw, &template.DeletedAt)

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

func (s *TemplatesStorage) GetAllMonthly(page int, limit int) ([]models.Template, error) {
	results, err := s.db.Query("select * from \"app\".\"templates\" where scheduling ->> 'type' = 'monthly' and \"deleted_at\" is null limit $1 offset $2", limit, (page-1)*limit)

	if err != nil {
		return nil, err
	}
	defer results.Close()

	templates := []models.Template{}

	for results.Next() {
		var template models.Template
		var schedulingRaw []byte

		err = results.Scan(&template.ID, &template.Name, &template.Description, &template.StartsAt, &template.EndsAt, &template.Enabled, &schedulingRaw, &template.DeletedAt)

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
