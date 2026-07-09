package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/chop1k/medods-test/internal/handlers"
)

// RegisterTemplateRoutes wires up the /tasks/templates resource routes
// described in the OpenAPI spec under the "Tasks" tag:
//
//	GET    /tasks/templates
//	POST   /tasks/templates
//	GET    /tasks/templates/{template_id}
//	PUT    /tasks/templates/{template_id}
//	DELETE /tasks/templates/{template_id}
func RegisterTemplateRoutes(router *gin.RouterGroup, h *handlers.TemplateHandler) {
	templates := router.Group("/tasks/templates")
	{
		templates.GET("", h.GetTemplates)
		templates.POST("", h.CreateTemplate)
		templates.GET("/:template_id", h.GetTemplateByID)
		templates.PUT("/:template_id", h.UpdateTemplate)
		templates.DELETE("/:template_id", h.DeleteTemplate)
	}
}

// RegisterTaskRoutes wires up the /tasks resource routes described in the
// OpenAPI spec under the "Tasks" tag:
//
//	GET    /tasks
//	POST   /tasks
//	GET    /tasks/{task_id}
//	PUT    /tasks/{task_id}
//	DELETE /tasks/{task_id}
func RegisterTaskRoutes(router *gin.RouterGroup, h *handlers.TaskHandler) {
	tasks := router.Group("/tasks")
	{
		tasks.GET("", h.GetTasks)
		tasks.POST("", h.CreateTask)
		tasks.GET("/:task_id", h.GetTaskByID)
		tasks.PUT("/:task_id", h.UpdateTask)
		tasks.DELETE("/:task_id", h.DeleteTask)
	}
}

// RegisterTagRoutes wires up the /tasks/tags resource routes described in
// the OpenAPI spec under the "Grouping" tag:
//
//	GET    /tasks/tags
//	POST   /tasks/tags
//	GET    /tasks/tags/{tag_id}
//	PUT    /tasks/tags/{tag_id}
//	DELETE /tasks/tags/{tag_id}
func RegisterTagRoutes(router *gin.RouterGroup, h *handlers.TagHandler) {
	tags := router.Group("/tasks/tags")
	{
		tags.GET("", h.GetTags)
		tags.POST("", h.CreateTag)
		tags.GET("/:tag_id", h.GetTagByID)
		tags.PUT("/:tag_id", h.UpdateTag)
		tags.DELETE("/:tag_id", h.DeleteTag)
	}
}

func RegisterSchedulingRoutes(router *gin.RouterGroup, h *handlers.SchedulingHandler) {
	tags := router.Group("/scheduling")
	{
		tags.POST("/connectivity-test", h.ConnectivityTest)
		tags.POST("/daily-cron-hook", h.DailyCronHook)
		tags.GET("/calendar", h.GetCalendar)
	}
}
