package service

import (
	"TaskSync/internal/entities"
	"TaskSync/internal/storage"
	"context"
	"time"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type People interface {
	Create(ctx context.Context, people entities.People) (int, error)
	GetByID(ctx context.Context, peopleID int) (entities.People, error)
	GetByFilter(ctx context.Context, filterPeople entities.People, limit, offset int) ([]entities.People, error)
	List(ctx context.Context) ([]entities.People, error)
	Update(ctx context.Context, people entities.People) error
	Delete(ctx context.Context, peopleID int) error
}

type Task interface {
	Create(ctx context.Context, peopleID int, task entities.Task) (int, error)
	GetByID(ctx context.Context, taskID int) (entities.Task, error)
	List(ctx context.Context) ([]entities.Task, error)
	Update(ctx context.Context, task entities.Task) error
	UpdatePeople(ctx context.Context, peopleID, taskID int) error
	Delete(ctx context.Context, taskID int) error
}

// управление временем выполнения
type Time interface {
	StartTimeEntry(ctx context.Context, task entities.Task) error
	EndTimeEntry(ctx context.Context, task entities.Task) error
	GetTaskTimeSpent(ctx context.Context, peopleID int, startTime, endTime time.Time) ([]entities.TaskTimeSpent, error)
}

type Service struct {
	People
	Task
	Time
}

func NewService(s *storage.Storage) *Service {
	return &Service{
		People: NewPeopleService(s.PeopleManage),
		Task:   NewTaskService(s.TaskManage),
		Time:   NewTimeService(s.TimeManage),
	}
}
