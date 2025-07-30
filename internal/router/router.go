package router

import (
	"task_service/internal/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router(taskHandler *handler.TaskHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/tasks", taskHandler.CreateTask)
	r.GET("/tasks/:id", taskHandler.GetTaskStatus)
	r.POST("/tasks/:id/files", taskHandler.AddFileToTask)
	r.GET("/archive/:id", taskHandler.DownloadArchive)

	return r
}
