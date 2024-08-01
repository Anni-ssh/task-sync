package storage

import (
	"TaskSync/internal/entities"
	"TaskSync/internal/storage/postgres"
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

//go:generate mockgen -source=storage.go -destination=mocks/mock.go

type PeopleManage interface {
	Create(ctx context.Context, people entities.People) (int, error)
	GetByID(ctx context.Context, peopleID int) (entities.People, error)
	GetByFilter(ctx context.Context, filterPeople entities.People, limit, offset int) ([]entities.People, error)
	List(ctx context.Context) ([]entities.People, error)
	Update(ctx context.Context, people entities.People) error
	Delete(ctx context.Context, peopleID int) error
}

type TaskManage interface {
	Create(ctx context.Context, task entities.Task) (int, error)
	GetByID(ctx context.Context, taskID int) (entities.Task, error)
	List(ctx context.Context) ([]entities.Task, error)
	Update(ctx context.Context, taskID int, title string, description string) error
	UpdatePeople(ctx context.Context, peopleID, taskID int) error
	Delete(ctx context.Context, taskID int) error
}

// управление временем выполнения
type TimeManage interface {
	StartTimeEntry(ctx context.Context, taskID int, timeEntries time.Time) error
	EndTimeEntry(ctx context.Context, taskID int, endTime time.Time) error
	GetTaskTimeSpent(ctx context.Context, peopleID int, startTime, endTime time.Time) ([]entities.TaskTimeSpent, error)
}

type Storage struct {
	PeopleManage
	TaskManage
	TimeManage
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		PeopleManage: postgres.NewPeopleManage(db),
		TaskManage:   postgres.NewTaskManage(db),
		TimeManage:   postgres.NewTimeManage(db),
	}
}
