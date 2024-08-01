package postgres

import (
	"TaskSync/internal/entities"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type TaskManagePostgres struct {
	db *sql.DB
}

func NewTaskManage(db *sql.DB) *TaskManagePostgres {
	return &TaskManagePostgres{db: db}
}

func (t *TaskManagePostgres) Create(ctx context.Context, task entities.Task) (int, error) {

	const op = "postgres.Task.Create"

	tx, err := t.db.BeginTx(ctx, nil)

	if err != nil {
		return 0, fmt.Errorf("database error: %w, operation: %s", err, op)
	}

	q := `INSERT INTO tasks (title, description) 
      VALUES($1, $2)
	  RETURNING id;`

	var newTaskID int

	err = tx.QueryRowContext(ctx, q, task.Title, task.Description).Scan(&newTaskID)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("database error: %w, operation: %s", err, op)
	}

	q = `INSERT INTO time_entries (people_id, task_id, start_time, end_time) 
      VALUES($1, $2, $3, $4)
	  RETURNING id;`

	result, err := tx.ExecContext(ctx, q, task.TimeEntry.PeopleID, newTaskID, task.TimeEntry.StartTime, task.TimeEntry.EndTime)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("database error: %w, operation: %s", err, op)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error retrieving affected rows: %w, operation: %s", err, op)
	}

	if rowsAffected == 0 {
		return 0, fmt.Errorf("no rows affected, operation: %s", op)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("database error: %w, operation: %s", err, op)
	}

	return newTaskID, nil

}

func (t *TaskManagePostgres) GetByID(ctx context.Context, taskID int) (entities.Task, error) {
	const op = "postgres.Task.GetByID"

	q := `SELECT t.id, t.title, t.description, te.people_id, te.start_time, te.end_time, te.created_at 
	FROM tasks t
	JOIN time_entries te ON t.id = te.task_id
	WHERE t.id = $1;`

	var task entities.Task
	row := t.db.QueryRowContext(ctx, q, taskID)

	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.TimeEntry.PeopleID, &task.TimeEntry.StartTime, &task.TimeEntry.EndTime, &task.TimeEntry.Created)
	if err != nil {
		if err == sql.ErrNoRows {
			return task, fmt.Errorf("no records found, operation: %s", op)
		}
		return task, fmt.Errorf("scan error: %w, operation: %s", err, op)
	}

	return task, nil
}

func (t *TaskManagePostgres) Delete(ctx context.Context, taskID int) error {
	const op = "postgres.Task.Delete"

	// Удаление каскадное, вместе удаляется вся строка из таблицы time_entries
	q := `DELETE FROM tasks WHERE id = $1;`

	result, err := t.db.ExecContext(ctx, q, taskID)
	if err != nil {
		return fmt.Errorf("database error: %w, operation: %s", err, op)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error retrieving affected rows: %w, operation: %s", err, op)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected, operation: %s", op)
	}

	return nil
}

func (t *TaskManagePostgres) List(ctx context.Context) ([]entities.Task, error) {
	const op = "postgres.Task.List"

	q := `SELECT t.id, t.title, t.description, te.people_id, te.start_time, te.end_time, te.created_at 
	FROM tasks t
	JOIN time_entries te ON t.id = te.task_id;`

	var taskList []entities.Task

	rows, err := t.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("database error: %w, operation: %s", err, op)
	}

	for rows.Next() {
		var task entities.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.TimeEntry.PeopleID, &task.TimeEntry.StartTime, &task.TimeEntry.EndTime, &task.TimeEntry.Created)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w, operation: %s", err, op)
		}

		taskList = append(taskList, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w, operation: %s", err, op)
	}

	return taskList, nil
}

func (t *TaskManagePostgres) Update(ctx context.Context, taskID int, title string, description string) error {
	const op = "postgres.task.Update"

	// Конструктор строки для запроса
	var q strings.Builder
	q.WriteString(`UPDATE tasks SET`)

	var args []interface{}
	argCount := 1

	// Добавление значений в запрос
	if title != "" {
		q.WriteString(fmt.Sprintf(" title = $%d,", argCount))
		args = append(args, title)
		argCount++
	}

	if description != "" {
		q.WriteString(fmt.Sprintf(" description = $%d", argCount))
		args = append(args, description)
		argCount++
	}

	// Доабавление ID обновляемой записи
	q.WriteString(fmt.Sprintf(" WHERE id = $%d", argCount))
	args = append(args, taskID)

	result, err := t.db.ExecContext(ctx, q.String(), args...)
	if err != nil {
		return fmt.Errorf("database error: %w, operation: %s", err, op)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error retrieving affected rows: %w, operation: %s", err, op)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected, operation: %s", op)
	}

	return nil
}

func (t *TaskManagePostgres) UpdatePeople(ctx context.Context, peopleID, taskID int) error {
	const op = "postgres.Task.UpdatePeople"

	// Проверяем значения
	if peopleID <= 0 || taskID <= 0 {
		return fmt.Errorf("incorrect values or their absence, operation: %s", op)
	}

	q := `UPDATE time_entries 
		SET people_id = $1
		WHERE task_id = $2`

	result, err := t.db.ExecContext(ctx, q, peopleID, taskID)
	if err != nil {
		return fmt.Errorf("database error: %w, operation: %s", err, op)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error retrieving affected rows: %w, operation: %s", err, op)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected, operation: %s", op)
	}

	return nil
}
