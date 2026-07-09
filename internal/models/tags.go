package models

import "time"

type TagSortableField string

const (
	TagSortByID   TagSortableField = "id"
	TagSortByName TagSortableField = "name"
	TagSortByType TagSortableField = "type"
)

type TagType string

const (
	TagTypePredefined  TagType = "predefined"
	TagTypeUserDefined TagType = "user-defined"
)

type TagBody struct {
	Name      string     `json:"name" binding:"required,min=2,max=64"`
	Type      TagType    `json:"type" binding:"required,oneof=predefined user-defined"`
	DeletedAt *time.Time `json:"deleted_at" binding:"omitempty"`
}

type Tag struct {
	ID int `json:"id"`
	TagBody
}

type TagListResponse struct {
	Data []Tag          `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type ListTagsQuery struct {
	Page      int    `form:"page,default=1" binding:"min=1"`
	Limit     int    `form:"limit,default=20" binding:"min=1,max=100"`
	Sort      string `form:"sort,default=asc" binding:"omitempty,oneof=asc desc"`
	SortField string `form:"sort-field,default=id" binding:"omitempty,oneof=id name type"`
}

type TagIDParam struct {
	TagID int `uri:"tag_id" binding:"required,min=1"`
}
