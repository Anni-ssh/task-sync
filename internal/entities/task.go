package entities

import (
	"time"
)

// Структура для хранения данных о времени
type TimeEntry struct {
	ID        int       `json:"id"`
	PeopleID  int       `json:"people_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Created   time.Time `json:"created"`
}

// Duration затраченное время на выполнение задачи, разница между StartTime и EndTime.
func (t *TimeEntry) Duration() time.Duration {
	return t.EndTime.Sub(t.StartTime)
}

// Структура для задачи
type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	TimeEntry   TimeEntry
}

// Структура для вывода трудозатрат по пользователю определённый период.
type TaskTimeSpent struct {
	PeopleID   int    `json:"people_id"`
	Surname    string `json:"surname"`
	Name       string `json:"name"`
	Patronymic string `json:"patronymic"`
	TaskID     int    `json:"task_id"`
	TaskTitle  string `json:"task_title"`
	TimeSpent  string `json:"time_spent"`
}
