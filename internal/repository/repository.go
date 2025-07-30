package repository

import (
	"errors"
	"sync"
	"task_service/internal/models"

	"github.com/google/uuid"
)

var (
	ErrTaskNotFound   = errors.New("задача не найдена")
	ErrTaskFilesLimit = errors.New("достигнут лимит файлов задачи")
)

type TaskRepository interface {
	CreateTask() (models.Task, error)
	AddFile(taskID uuid.UUID, file models.FileInfo) error
	GetTask(taskID uuid.UUID) (models.Task, error)
	ListActiveTasks() []models.Task
	SetStatus(taskID uuid.UUID, status models.TaskStatus) error
	SetArchiveURL(taskID uuid.UUID, url string) error
	AddError(taskID uuid.UUID, err models.FileError) error
}

type TaskRepo struct {
	mu    sync.RWMutex
	tasks map[uuid.UUID]*models.Task
}

func NewTaskRepo() *TaskRepo {
	return &TaskRepo{
		tasks: make(map[uuid.UUID]*models.Task),
	}
}

func (r *TaskRepo) CreateTask() (models.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	task := &models.Task{
		ID:     uuid.New(),
		Status: models.StatusCreated,
		Files:  make([]models.FileInfo, 0, 3),
	}

	r.tasks[task.ID] = task
	return *task, nil
}

func (r *TaskRepo) AddFile(taskID uuid.UUID, file models.FileInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	task, ok := r.tasks[taskID]
	if !ok {
		return ErrTaskNotFound
	}
	if len(task.Files) >= 3 {
		return ErrTaskFilesLimit
	}
	task.Files = append(task.Files, file)
	return nil
}

func (r *TaskRepo) GetTask(taskID uuid.UUID) (models.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	task, ok := r.tasks[taskID]
	if !ok {
		return models.Task{}, ErrTaskNotFound
	}
	return *task, nil
}

func (r *TaskRepo) ListActiveTasks() []models.Task {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tasks := make([]models.Task, 0)
	for _, t := range r.tasks {
		if t.Status == models.StatusCreated || t.Status == models.StatusActive {
			tasks = append(tasks, *t)
		}
	}
	return tasks
}

func (r *TaskRepo) SetStatus(taskID uuid.UUID, status models.TaskStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	task, ok := r.tasks[taskID]
	if !ok {
		return ErrTaskNotFound
	}
	task.Status = status
	return nil
}

func (r *TaskRepo) SetArchiveURL(taskID uuid.UUID, url string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	task, ok := r.tasks[taskID]
	if !ok {
		return ErrTaskNotFound
	}
	task.ArchiveURL = url
	return nil
}

func (r *TaskRepo) AddError(taskID uuid.UUID, err models.FileError) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	task, ok := r.tasks[taskID]
	if !ok {
		return ErrTaskNotFound
	}
	task.Errors = append(task.Errors, err)
	return nil
}
