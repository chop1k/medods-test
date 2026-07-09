package models

import "time"

// TagSortableField enumerates the fields tags can be sorted by, mirroring
// the `TagSortableField` schema.
type TagSortableField string

const (
	TagSortByID   TagSortableField = "id"
	TagSortByName TagSortableField = "name"
	TagSortByType TagSortableField = "type"
)

// TagType enumerates the `type` field values on TagBody.
type TagType string

const (
	TagTypePredefined  TagType = "predefined"
	TagTypeUserDefined TagType = "user-defined"
)

// TagBody mirrors the `TagBody` schema and is used as the request payload
// for create/update operations.
//
// NOTE: the spec's PUT /tasks/tags/{tag_id} request body actually references
// the full `Tag` schema rather than `TagBody`, but - consistent with how
// TemplateBody/TaskBody are reused for updates elsewhere in this codebase -
// TagBody is used here too; the path's `tag_id` is the source of truth for
// the identifier.
type TagBody struct {
	Name      string     `json:"name" binding:"required,min=2,max=64"`
	Type      TagType    `json:"type" binding:"required,oneof=predefined user-defined"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// Tag mirrors the `Tag` schema (IdentifiableResource + TagBody).
type Tag struct {
	ID int64 `json:"id"`
	TagBody
}

// TagListResponse is the envelope returned by GET /tasks/tags.
type TagListResponse struct {
	Data []Tag          `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// ListTagsQuery captures the query parameters accepted by GET /tasks/tags.
type ListTagsQuery struct {
	Page      int    `form:"page,default=1" binding:"min=1"`
	Limit     int    `form:"limit,default=20" binding:"min=1,max=100"`
	Sort      string `form:"sort,default=asc" binding:"omitempty,oneof=asc desc"`
	SortField string `form:"sort-field,default=id" binding:"omitempty,oneof=id name type"`
}

// TagIDParam captures the `tag_id` path parameter shared by
// GET/PUT/DELETE /tasks/tags/{tag_id}.
type TagIDParam struct {
	TagID int64 `uri:"tag_id" binding:"required,min=1"`
}
