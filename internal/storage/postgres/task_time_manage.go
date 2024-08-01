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

func (t *TimeManagePostgres) StartTimeEntry(ctx context.Context, taskID int, timeEntries time.Time) error {
	const op = "postgres.Time.StartTimeEntry"

	q := `UPDATE time_entries 
		SET start_time = $1
		WHERE task_id = $2;`

	result, err := t.db.ExecContext(ctx, q, timeEntries, taskID)
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
func (t *TimeManagePostgres) EndTimeEntry(ctx context.Context, taskID int, endTime time.Time) error {

	const op = "postgres.Time.EndTimeEntry"

	q := `UPDATE time_entries 
		SET end_time = $1
		WHERE task_id = $2;`

	result, err := t.db.ExecContext(ctx, q, endTime, taskID)
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

// GetTaskTimeSpent возвращает трудозатраты по пользователю за определённый период.
func (t *TimeManagePostgres) GetTaskTimeSpent(ctx context.Context, peopleID int, startTime, endTime time.Time) ([]entities.TaskTimeSpent, error) {
	const q = `
		SELECT
			p.id AS people_id,
			p.surname,
			p.name,
			p.patronymic,
			t.id AS task_id,
			t.title AS task_title,
			COALESCE(SUM(EXTRACT(EPOCH FROM (te.end_time - te.start_time)) / 3600), 0) * INTERVAL '1 hour' AS time_spent
		FROM
			tasks t
		JOIN
			time_entries te ON t.id = te.task_id
		JOIN
			people_info p ON te.people_id = p.id
		WHERE
			p.id = $1 AND
			te.start_time >= DATE_TRUNC('day', $2::TIMESTAMP) AND
			te.end_time <= DATE_TRUNC('day', $3::TIMESTAMP) + INTERVAL '1 day'
		GROUP BY
			p.id, p.surname, p.name, p.patronymic, t.id, t.title
		ORDER BY
			time_spent DESC;
	`

	rows, err := t.db.QueryContext(ctx, q, peopleID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var entries []entities.TaskTimeSpent
	for rows.Next() {
		var entry entities.TaskTimeSpent
		if err := rows.Scan(&entry.PeopleID, &entry.Surname, &entry.Name, &entry.Patronymic, &entry.TaskID, &entry.TaskTitle, &entry.TimeSpent); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return entries, nil
}
