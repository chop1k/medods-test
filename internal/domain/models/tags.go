package models

type TagSortableField string

const (
	TagSortByID          TagSortableField = "id"
	TagSortByName        TagSortableField = "name"
	TagSortByDescription TagSortableField = "description"
)

type TagType string

const (
	TagTypePredefined  TagType = "predefined"
	TagTypeUserDefined TagType = "user-defined"
)

type TagBody struct {
	Name        string   `json:"name" binding:"required,min=2,max=64"`
	Description *string  `json:"description,omitempty" binding:"omitempty,min=2"`
	Type        *TagType `json:"type,omitempty"`
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
