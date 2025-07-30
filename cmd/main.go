package main

import (
	_ "task_service/docs"
	"task_service/internal/handler"
	"task_service/internal/logger"
	"task_service/internal/repository"
	"task_service/internal/router"
	"task_service/internal/service"

	"go.uber.org/zap"
)

// @Title  Task Service
func main() {
	if err := logger.Init(); err != nil {
		panic(err)
	}

	defer logger.Sync()

	log := logger.L()

	repo := repository.NewTaskRepo()
	service := service.NewTaskService(repo, log)
	taskHandler := handler.NewTaskHandler(service)
	r := router.Router(taskHandler)
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Не удалось запустить сервер: ", zap.Error(err))
	}

}
