package service

import (
	"TaskSync/internal/entities"
	"TaskSync/internal/storage"
	"context"
)

// TaskService представляет сервис для работы с данными задач.
type TaskService struct {
	storage storage.TaskManage
}

// NewTaskService создает новый экземпляр TaskService.
func NewTaskService(t storage.TaskManage) *TaskService {
	return &TaskService{storage: t}
}

// Create создает новую задачу для пользователя.
func (t *TaskService) Create(ctx context.Context, task entities.Task) (int, error) {
	return t.storage.Create(ctx, task)
}

// GetByID возвращает данные задачи по её ID.
func (t *TaskService) GetByID(ctx context.Context, taskID int) (entities.Task, error) {
	return t.storage.GetByID(ctx, taskID)
}

// List возвращает список всех задач.
func (t *TaskService) List(ctx context.Context) ([]entities.Task, error) {
	return t.storage.List(ctx)
}

// Update обновляет данные задачи.
func (t *TaskService) Update(ctx context.Context, taskID int, title string, description string) error {
	return t.storage.Update(ctx, taskID, title, description)
}

// UpdatePeople обновляет исполнителя задачи.
func (t *TaskService) UpdatePeople(ctx context.Context, peopleID, taskID int) error {
	return t.storage.UpdatePeople(ctx, peopleID, taskID)
}

// Delete удаляет задачу по её ID.
func (t *TaskService) Delete(ctx context.Context, taskID int) error {
	return t.storage.Delete(ctx, taskID)
}
