package repository

import (
	"task_service/internal/models"
	"testing"

	"github.com/google/uuid"
)

func TestCreateTask(t *testing.T) {
	repo := NewTaskRepo()
	task, err := repo.CreateTask()
	if err != nil {
		t.Fatalf("CreateTask error: %v", err)
	}
	if task.Status != models.StatusCreated {
		t.Errorf("Expected status 'created', got %v", task.Status)
	}
	if len(task.Files) != 0 {
		t.Errorf("Expected empty files, got %d", len(task.Files))
	}
}

func TestAddFileAndLimit(t *testing.T) {
	repo := NewTaskRepo()
	task, _ := repo.CreateTask()
	file := models.FileInfo{URL: "url", Name: "file1.pdf", Type: ".pdf"}
	for i := 0; i < 3; i++ {
		err := repo.AddFile(task.ID, file)
		if err != nil {
			t.Fatalf("AddFile error: %v", err)
		}
	}
	// Превышение лимита
	err := repo.AddFile(task.ID, file)
	if err != ErrTaskFilesLimit {
		t.Errorf("Expected ErrTaskFilesLimit, got %v", err)
	}
}

func TestSetStatusAndArchiveURL(t *testing.T) {
	repo := NewTaskRepo()
	task, _ := repo.CreateTask()
	err := repo.SetStatus(task.ID, models.StatusActive)
	if err != nil {
		t.Fatalf("SetStatus error: %v", err)
	}
	task2, _ := repo.GetTask(task.ID)
	if task2.Status != models.StatusActive {
		t.Errorf("Expected status 'active', got %v", task2.Status)
	}
	err = repo.SetArchiveURL(task.ID, "archive.zip")
	if err != nil {
		t.Fatalf("SetArchiveURL error: %v", err)
	}
	task2, _ = repo.GetTask(task.ID)
	if task2.ArchiveURL != "archive.zip" {
		t.Errorf("Expected archive.zip, got %v", task2.ArchiveURL)
	}
}

func TestAddError(t *testing.T) {
	repo := NewTaskRepo()
	task, _ := repo.CreateTask()
	err := repo.AddError(task.ID, models.FileError{URL: "bad", Reason: "fail"})
	if err != nil {
		t.Fatalf("AddError error: %v", err)
	}
	task2, _ := repo.GetTask(task.ID)
	if len(task2.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(task2.Errors))
	}
}

func TestListActiveTasks(t *testing.T) {
	repo := NewTaskRepo()
	for i := 0; i < 2; i++ {
		task, _ := repo.CreateTask()
		repo.SetStatus(task.ID, models.StatusCreated)
	}
	tasks := repo.ListActiveTasks()
	if len(tasks) != 2 {
		t.Errorf("Expected 2 active tasks, got %d", len(tasks))
	}
}

func TestGetTaskNotFound(t *testing.T) {
	repo := NewTaskRepo()
	_, err := repo.GetTask(uuid.New())
	if err != ErrTaskNotFound {
		t.Errorf("Expected ErrTaskNotFound, got %v", err)
	}
}
