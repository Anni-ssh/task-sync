package postgres

import (
	"TaskSync/internal/entities"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type TimeManagePostgres struct {
	db *sql.DB
}

func NewTimeManage(db *sql.DB) *TimeManagePostgres {
	return &TimeManagePostgres{db: db}
}

func (t *TimeManagePostgres) StartTimeEntry(ctx context.Context, taskID int, startTime time.Time) error {
	const op = "postgres.Time.StartTimeEntry"

	query := `UPDATE time_entries 
		SET start_time = $1
		WHERE task_id = $2;`

	stmt, err := t.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare error: %w, operation: %s", err, op)
	}

	result, err := stmt.ExecContext(ctx, startTime, taskID)
	if err != nil {
		return fmt.Errorf("failed to update start time for task ID %d: %w, operation: %s", taskID, err, op)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error retrieving affected rows: %w, operation: %s", err, op)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows updated for task ID %d, operation: %s", taskID, op)
	}

	return nil
}

func (t *TimeManagePostgres) EndTimeEntry(ctx context.Context, taskID int, endTime time.Time) error {
	const op = "postgres.Time.EndTimeEntry"

	query := `UPDATE time_entries 
		SET end_time = $1
		WHERE task_id = $2;`

	stmt, err := t.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare error: %w, operation: %s", err, op)
	}

	result, err := stmt.ExecContext(ctx, endTime, taskID)
	if err != nil {
		return fmt.Errorf("failed to update end time for task ID %d: %w, operation: %s", taskID, err, op)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error retrieving affected rows: %w, operation: %s", err, op)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows updated for task ID %d, operation: %s", taskID, op)
	}

	return nil
}

// GetTaskTimeSpent извлекает данные о том, сколько времени пользователь потратил на задачи за определённый период времени.
// Функция возвращает список, в котором содержится информация о пользователе, задачах и количестве времени, затраченного на каждую задачу.
func (t *TimeManagePostgres) TasksTimeSpent(ctx context.Context, peopleID int, startTime, endTime time.Time) ([]entities.TaskTimeSpent, error) {
	const op = "postgres.Time.GetTaskTimeSpent"

	// Определяем запрос
	const query = `
	SELECT
		p.id AS people_id,
		p.surname,
		p.name,
		p.patronymic,
		t.id AS task_id,
		t.title AS task_title,
    COALESCE(
        SUM(
            EXTRACT(EPOCH FROM (te.end_time - te.start_time)) / 3600
        ),
        0
    ) * INTERVAL '1 hour' AS time_spent
	FROM
		tasks t
	JOIN
		time_entries te ON t.id = te.task_id
	JOIN
		people_info p ON te.people_id = p.id
	WHERE
		p.id = $1
		AND te.start_time >= $2::timestamptz
		AND te.end_time <= $3::timestamptz
	GROUP BY
		p.id, p.surname, p.name, p.patronymic, t.id, t.title
	ORDER BY
		time_spent DESC;

	`

	// Подготовка запроса
	stmt, err := t.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("prepare error: %w, operation: %s", err, op)
	}

	// Выполнение подготовленного запроса
	rows, err := stmt.QueryContext(ctx, peopleID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("query error: %w, operation: %s", err, op)
	}

	// Получение результатов
	var entries []entities.TaskTimeSpent

	for rows.Next() {
		var entry entities.TaskTimeSpent
		if err := rows.Scan(&entry.PeopleID, &entry.Surname, &entry.Name, &entry.Patronymic, &entry.TaskID, &entry.TaskTitle, &entry.TimeSpent); err != nil {
			return nil, fmt.Errorf("scan error: %w, operation: %s", err, op)
		}
		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w, operation: %s", err, op)
	}

	return entries, nil
}
