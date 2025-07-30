package service

import (
	"errors"
	"os"
	"task_service/internal/models"
	"task_service/internal/repository"
	"testing"

	"go.uber.org/zap"
)

func TestCreateTaskLimit(t *testing.T) {
	repo := repository.NewTaskRepo()
	log := zap.NewNop()
	service := NewTaskService(repo, log)
	for i := 0; i < 3; i++ {
		_, err := service.CreateTask()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	_, err := service.CreateTask()
	if err == nil {
		t.Errorf("expected error on 4th task, got nil")
	}
}

func TestAddFileExtension(t *testing.T) {
	repo := repository.NewTaskRepo()
	log := zap.NewNop()
	service := NewTaskService(repo, log)
	task, _ := service.CreateTask()
	// Недопустимое расширение
	err := service.AddFile(task.ID, "http://site/file.txt")
	if err == nil || err.Error() != "недопустимое расширение файла" {
		t.Errorf("expected extension error, got %v", err)
	}
}

func TestAddFileLimit(t *testing.T) {
	repo := repository.NewTaskRepo()
	log := zap.NewNop()
	service := NewTaskService(repo, log)
	task, _ := service.CreateTask()
	for i := 0; i < 3; i++ {
		// Мокаем скачивание
		service.repo.AddFile(task.ID, models.FileInfo{URL: "url", Name: "file.pdf", Type: ".pdf"})
	}
	err := service.AddFile(task.ID, "http://site/file.pdf")
	if err == nil {
		t.Errorf("expected error on 4th file, got nil")
	}
}

func TestDeleteArchiveRemovesFiles(t *testing.T) {
	repo := repository.NewTaskRepo()
	log := zap.NewNop()
	service := NewTaskService(repo, log)
	task, _ := service.CreateTask()
	// Создаём временные файлы

	tempDir := "temp"
	os.MkdirAll(tempDir, 0755)
	f1, _ := os.CreateTemp(tempDir, "testfile1")
	f2, _ := os.CreateTemp(tempDir, "testfile2")
	f1.Close()
	f2.Close()
	repo.AddFile(task.ID, models.FileInfo{URL: "url1", Name: f1.Name(), Type: ".pdf"})
	repo.AddFile(task.ID, models.FileInfo{URL: "url2", Name: f2.Name(), Type: ".pdf"})
	// Удаляем
	err := service.DeleteArchive(task.ID)
	if err != nil {
		t.Fatalf("DeleteArchive error: %v", err)
	}
	if _, err := os.Stat(f1.Name()); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("file1 not deleted")
	}
	if _, err := os.Stat(f2.Name()); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("file2 not deleted")
	}
}
