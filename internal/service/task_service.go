package service

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"task_service/internal/models"
	"task_service/internal/repository"
	"task_service/internal/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TaskService struct {
	repo repository.TaskRepository
	log  *zap.Logger
}

func NewTaskService(repo repository.TaskRepository, log *zap.Logger) *TaskService {
	return &TaskService{
		repo: repo,
		log:  log,
	}
}

func (s *TaskService) CreateTask() (models.Task, error) {
	if len(s.repo.ListActiveTasks()) >= 3 {
		s.log.Warn("достигнут лимит активных задач")
		return models.Task{}, errors.New("достигнут лимит активных задач")
	}
	s.log.Info("Создание новой задачи")
	return s.repo.CreateTask()
}

func (s *TaskService) DeleteArchive(taskID uuid.UUID) error {
	task, err := s.repo.GetTask(taskID)
	if err != nil {
		return err
	}
	for _, f := range task.Files {
		if f.Name != "" {
			if err := os.Remove(f.Name); err != nil {
				s.log.Warn("Не удалось удалить файл задачи", zap.String("file", f.Name), zap.Error(err))
			}
		}
	}
	task.Files = nil
	return nil
}

func (s *TaskService) AddFile(taskID uuid.UUID, url string) error {
	if !utils.IsAllowedExtension(url) {
		s.log.Warn("недопустимое расширение файла")
		return errors.New("недопустимое расширение файла")
	}
	localPath, err := utils.DownloadFile(url)
	if err != nil {
		ferr := models.FileError{URL: url, Reason: err.Error()}
		_ = s.repo.AddError(taskID, ferr)
		return fmt.Errorf("ошибка скачивания файла: %w", err)
	}
	file := models.FileInfo{
		URL:        url,
		Name:       localPath,
		Type:       filepathExt(url),
		Downloader: true,
	}
	err = s.repo.AddFile(taskID, file)
	if err != nil {
		s.log.Error("ошибка добавления файла в задачу", zap.Error(err))
		return err
	}
	task, _ := s.repo.GetTask(taskID)
	if len(task.Files) == 3 {
		s.log.Info("Все файлы загружены, начинаем архивирование")
		s.repo.SetStatus(taskID, models.StatusActive)
		archiveName := fmt.Sprintf("archive_%s.zip", taskID.String())
		filePaths := make([]string, 0, 3)
		for _, f := range task.Files {
			filePaths = append(filePaths, f.Name)
		}
		archivePath, err := utils.CreateZipArchive(filePaths, archiveName)
		if err != nil {
			s.log.Error("ошибка создания архива", zap.Error(err))
			s.repo.SetStatus(taskID, models.StatusError)
			s.repo.AddError(taskID, models.FileError{URL: "archive", Reason: err.Error()})
			return fmt.Errorf("ошибка создания архива: %w", err)
		}
		s.repo.SetArchiveURL(taskID, archivePath)
		s.repo.SetStatus(taskID, models.StatusDone)
	}
	s.repo.SetStatus(taskID, models.StatusActive)
	return nil
}

func (s *TaskService) GetTaskStatus(taskID uuid.UUID) (models.Task, error) {
	s.log.Info("Получение статуса задачи", zap.String("taskID", taskID.String()))
	return s.repo.GetTask(taskID)
}

func filepathExt(url string) string {
	parts := strings.Split(url, ".")
	if len(parts) > 1 {
		return "." + parts[len(parts)-1]
	}
	return ""
}
