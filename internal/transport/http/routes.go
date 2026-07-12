package http

import (
	"github.com/chop1k/medods-test/internal/transport/http/handler"
	"github.com/gin-gonic/gin"
)

func RegisterTemplateRoutes(router *gin.RouterGroup, h *handler.TemplateHandler) {
	templates := router.Group("/tasks/templates")
	{
		templates.GET("", h.GetTemplates)
		templates.POST("", h.CreateTemplate)
		templates.GET("/:template_id", h.GetTemplateByID)
		templates.PUT("/:template_id", h.UpdateTemplate)
		templates.DELETE("/:template_id", h.DeleteTemplate)
	}
}

func RegisterTaskRoutes(router *gin.RouterGroup, h *handler.TaskHandler) {
	tasks := router.Group("/tasks")
	{
		tasks.GET("", h.GetTasks)
		tasks.POST("", h.CreateTask)
		tasks.GET("/:task_id", h.GetTaskByID)
		tasks.PUT("/:task_id", h.UpdateTask)
		tasks.DELETE("/:task_id", h.DeleteTask)
	}
}

func RegisterTagRoutes(router *gin.RouterGroup, h *handler.TagHandler) {
	tags := router.Group("/grouping/tags")
	{
		tags.GET("", h.GetTags)
		tags.POST("", h.CreateTag)
		tags.GET("/:tag_id", h.GetTagByID)
		tags.DELETE("/:tag_id", h.DeleteTag)
	}
}

func RegisterSchedulingRoutes(router *gin.RouterGroup, h *handler.SchedulingHandler) {
	scheduling := router.Group("/scheduling")
	{
		scheduling.POST("/connectivity-test", h.ConnectivityTest)
		scheduling.POST("/daily-cron-hook", h.DailyCronHook)
		scheduling.GET("/calendar", h.GetCalendar)
	}
}

// func RegisterSwaggerRoutes(router *gin.RouterGroup, h *handlers.SchedulingHandler) {
// 	swagger := v5emb.New()

// 	docs := router.Group("/docs")
// 	{
// 		docs.GET("", gin.WrapH(swagger))
// 	}
// }
