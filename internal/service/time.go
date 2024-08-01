package service

import (
	"TaskSync/internal/entities"
	"TaskSync/internal/storage"
	"context"
	"time"
)

// TimeService представляет сервис для работы с данными времени задач.
type TimeService struct {
	storage storage.TimeManage
}

// NewTimeService создает новый экземпляр TimeService.
func NewTimeService(t storage.TimeManage) *TimeService {
	return &TimeService{storage: t}
}

// StartTimeEntry начинает запись времени для задачи.
func (t *TimeService) StartTimeEntry(ctx context.Context, taskID int, timeEntries time.Time) error {
	return t.storage.StartTimeEntry(ctx, taskID, timeEntries)
}

// EndTimeEntry завершает запись времени для задачи.
func (t *TimeService) EndTimeEntry(ctx context.Context, task entities.Task) error {
	return t.storage.EndTimeEntry(ctx, task)
}

// GetTaskTimeSpent возвращает трудозатраты по пользователю за заданный период.
func (t *TimeService) GetTaskTimeSpent(ctx context.Context, peopleID int, startTime, endTime time.Time) ([]entities.TaskTimeSpent, error) {
	return t.storage.GetTaskTimeSpent(ctx, peopleID, startTime, endTime)
}
