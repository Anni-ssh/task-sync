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

	// Создание транзакции
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("database error: %w, operation: %s", err, op)
	}

	// Подготовка первого запроса
	insertTaskQuery := `INSERT INTO tasks (title, description) 
      VALUES($1, $2)
	  RETURNING id;`
	stmtInsertTask, err := tx.PrepareContext(ctx, insertTaskQuery)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("prepare error for insertTask: %w, operation: %s", err, op)
	}

	// Выполнение первого запроса
	var newTaskID int
	err = stmtInsertTask.QueryRowContext(ctx, task.Title, task.Description).Scan(&newTaskID)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("database error during insertTask execution: %w, operation: %s", err, op)
	}

	// Подготовка второго запроса
	insertTimeEntryQuery := `INSERT INTO time_entries (people_id, task_id, start_time, end_time) 
      VALUES($1, $2, $3, $4)
	  RETURNING id;`
	stmtInsertTimeEntry, err := tx.PrepareContext(ctx, insertTimeEntryQuery)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("prepare error for insertTimeEntry: %w, operation: %s", err, op)
	}

	// Выполнение второго запроса
	result, err := stmtInsertTimeEntry.ExecContext(ctx, task.TimeEntry.PeopleID, newTaskID, task.TimeEntry.StartTime, task.TimeEntry.EndTime)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("database error during insertTimeEntry execution: %w, operation: %s", err, op)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error retrieving affected rows: %w, operation: %s", err, op)
	}

	if rowsAffected == 0 {
		tx.Rollback()
		return 0, fmt.Errorf("no rows affected, operation: %s", op)
	}

	// Завершение транзакции
	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("database error during commit: %w, operation: %s", err, op)
	}

	return newTaskID, nil
}

func (t *TaskManagePostgres) GetByID(ctx context.Context, taskID int) (entities.Task, error) {
	const op = "postgres.Task.GetByID"

	query := `SELECT t.id, t.title, t.description, te.people_id, te.start_time, te.end_time, te.created_at 
	FROM tasks t
	JOIN time_entries te ON t.id = te.task_id
	WHERE t.id = $1;`

	stmt, err := t.db.PrepareContext(ctx, query)
	if err != nil {
		return entities.Task{}, fmt.Errorf("prepare error: %w, operation: %s", err, op)
	}

	var task entities.Task
	row := stmt.QueryRowContext(ctx, taskID)

	err = row.Scan(&task.ID, &task.Title, &task.Description, &task.TimeEntry.PeopleID, &task.TimeEntry.StartTime, &task.TimeEntry.EndTime, &task.TimeEntry.Created)
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

	query := `DELETE FROM tasks WHERE id = $1;`

	stmt, err := t.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare error: %w, operation: %s", err, op)
	}

	result, err := stmt.ExecContext(ctx, taskID)
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

	query := `SELECT t.id, t.title, t.description, te.people_id, te.start_time, te.end_time, te.created_at 
	FROM tasks t
	JOIN time_entries te ON t.id = te.task_id;`

	stmt, err := t.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("prepare error: %w, operation: %s", err, op)
	}

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("database error: %w, operation: %s", err, op)
	}

	var taskList []entities.Task

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

	var q strings.Builder
	q.WriteString(`UPDATE tasks SET`)

	var args []interface{}
	argCount := 1

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

	// Удаляем последнюю запятую, если она есть
	if q.Len() > len("UPDATE tasks SET") {
		query := q.String()
		query = query[:len(query)-1] // Удаление последней запятой
		q.Reset()
		q.WriteString(query)
	}

	// Добавление ID обновляемой записи
	q.WriteString(fmt.Sprintf(" WHERE id = $%d", argCount))
	args = append(args, taskID)

	query := q.String()

	stmt, err := t.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare error: %w, operation: %s", err, op)
	}

	result, err := stmt.ExecContext(ctx, args...)
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

	if peopleID <= 0 || taskID <= 0 {
		return fmt.Errorf("incorrect values or their absence, operation: %s", op)
	}

	query := `UPDATE time_entries 
		SET people_id = $1
		WHERE task_id = $2`

	stmt, err := t.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare error: %w, operation: %s", err, op)
	}

	result, err := stmt.ExecContext(ctx, peopleID, taskID)
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
