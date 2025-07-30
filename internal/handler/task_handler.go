package handler

import (
	"net/http"
	"strings"
	"task_service/internal/models"
	"task_service/internal/response"
	"task_service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TaskHandler отвечает за обработку задач
type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

// CreateTask godoc
// @Summary Создать новую задачу
// @Tags tasks
// @Accept json
// @Produce json
// @Success 201 {object} string "ID задачи"
// @Failure 429 {object} response.ErrorResponse "достигнут лимит активных задач"
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	task, err := h.service.CreateTask()
	if err != nil {
		c.JSON(http.StatusTooManyRequests, response.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": task.ID.String()})
}

// GetTaskStatus godoc
// @Summary Получить статус задачи
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "ID задачи"
// @Success 200 {object} response.TaskStatusResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTaskStatus(c *gin.Context) {
	id := c.Param("id")
	taskID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Некорректный формат ID"})
		return
	}
	task, err := h.service.GetTaskStatus(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.TaskStatusResponse{
		ID:     task.ID.String(),
		Status: string(task.Status),
	})
}

type AddFileRequest struct {
	URL string `json:"url" binding:"required"`
}

// AddFileToTask godoc
// @Summary Добавить файл в задачу
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "ID задачи"
// @Param file body AddFileRequest true "Ссылка на файл"
// @Success 200 {object} response.TaskStatusResponse "Статус задачи"
// @Failure 400 {object} response.ErrorResponse "Некорректный формат запроса"
// @Failure 404 {object} response.ErrorResponse "Задача не найдена"
// @Failure 422 {object} response.ErrorResponse "Недопустимое расширение файла"
// @Router /tasks/{id}/files [post]
func (h *TaskHandler) AddFileToTask(c *gin.Context) {
	id := c.Param("id")
	taskID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Некорректный формат ID"})
		return
	}
	var req AddFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Некорректный формат запроса"})
		return
	}
	err = h.service.AddFile(taskID, req.URL)
	if err != nil {
		if err.Error() == "недопустимое расширение файла" {
			c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, response.ErrorResponse{Error: err.Error()})
		return
	}
	task, _ := h.service.GetTaskStatus(taskID)
	archiveURL := strings.ReplaceAll(task.ArchiveURL, "\\", "/")
	c.JSON(http.StatusOK, response.TaskStatusResponse{
		ID:         task.ID.String(),
		Status:     string(task.Status),
		Files:      toFileInfoDTOs(task.Files),
		ArchiveURL: archiveURL,
		Errors:     toFileErrorDTOs(task.Errors),
	})
}

func toFileInfoDTOs(files []models.FileInfo) []response.FileInfoDTO {
	dtos := make([]response.FileInfoDTO, 0, len(files))
	for _, f := range files {
		dtos = append(dtos, response.FileInfoDTO{
			URL:  f.URL,
			Name: strings.ReplaceAll(f.Name, "\\", "/"),
			Type: f.Type,
		})
	}
	return dtos
}

func toFileErrorDTOs(errors []models.FileError) []response.FileErrorDTO {
	dtos := make([]response.FileErrorDTO, 0, len(errors))
	for _, e := range errors {
		dtos = append(dtos, response.FileErrorDTO{
			URL:    e.URL,
			Reason: e.Reason,
		})
	}
	return dtos
}

// DownloadArchive godoc
// @Summary Скачать архив задачи
// @Tags tasks
// @Accept json
// @Produce application/zip
// @Param id path string true "ID задачи"
// @Success 200 {file} file "ZIP архив"
// @Failure 404 {object} response.ErrorResponse "Задача или архив не найдены"
// @Router /archive/{id} [get]
func (h *TaskHandler) DownloadArchive(c *gin.Context) {
	id := c.Param("id")
	taskID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Некорректный формат ID"})
		return
	}
	task, err := h.service.GetTaskStatus(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{Error: err.Error()})
		return
	}
	if task.ArchiveURL == "" {
		c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Архив не найден"})
		return
	}
	archivePath := strings.ReplaceAll(task.ArchiveURL, "\\", "/")
	c.FileAttachment(archivePath, "archive.zip")
	// После отдачи — удалить архив и очистить ссылку
	_ = h.service.DeleteArchive(taskID)
}
