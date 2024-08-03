package postgres

import (
	"TaskSync/internal/entities"
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

type PeopleManagePostgres struct {
	db *sql.DB
}

func NewPeopleManage(db *sql.DB) *PeopleManagePostgres {
	return &PeopleManagePostgres{db: db}
}

func (p *PeopleManagePostgres) Create(ctx context.Context, people entities.People) (int, error) {
	const op = "postgres.People.Create"

	stmt, err := p.db.PrepareContext(ctx, `INSERT INTO people_info (passport_series, passport_number, surname, name, patronymic, address) 
	VALUES ($1, $2, $3, $4, $5, $6) 
	RETURNING id;`)
	if err != nil {
		return 0, fmt.Errorf("%s Prepare: %w", op, err)
	}

	var id int

	row := stmt.QueryRowContext(ctx, people.PassportSeries, people.PassportNumber, people.Surname, people.Name, people.Patronymic, people.Address)

	err = row.Scan(&id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "22023" { // "invalid_parameter_value"
				return 0, fmt.Errorf("%w, operation: %s", ErrInputData, op)
			}
			return 0, fmt.Errorf("database error: %w, operation: %s", pqErr, op)
		}
		return 0, fmt.Errorf("scan error: %w, operation: %s", err, op)
	}

	return id, nil
}

func (p *PeopleManagePostgres) GetByID(ctx context.Context, peopleID int) (entities.People, error) {
	const op = "postgres.People.Get"

	stmt, err := p.db.PrepareContext(ctx, `SELECT id, passport_series, passport_number, surname, name, patronymic, address FROM people_info WHERE id = $1;`)
	if err != nil {
		return entities.People{}, fmt.Errorf("prepare error: %w, operation: %s", err, op)
	}

	var people entities.People

	row := stmt.QueryRowContext(ctx, peopleID)

	err = row.Scan(&people.ID, &people.PassportSeries, &people.PassportNumber, &people.Surname, &people.Name, &people.Patronymic, &people.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			return people, fmt.Errorf("no records found, operation: %s", op)
		}
		return people, fmt.Errorf("scan error: %w, operation: %s", err, op)
	}

	return people, nil
}

func (p *PeopleManagePostgres) GetByFilter(ctx context.Context, filterPeople entities.People, limit, offset int) ([]entities.People, error) {
	const op = "postgres.People.GetByFilter"

	// Конструктор для запроса
	var q strings.Builder

	q.WriteString(`SELECT id, passport_series, passport_number, surname, name, patronymic, address 
	FROM people_info
	WHERE 1 = 1`)
	// При отсутствии фильтров - выведет все записи.

	argCount := 1

	// Собираем условия фильтрации
	var args []interface{}
	if filterPeople.ID != 0 {
		q.WriteString(fmt.Sprintf(" AND id = $%d", argCount))
		args = append(args, filterPeople.ID)
		argCount++
	}
	if filterPeople.PassportSeries != 0 {
		q.WriteString(fmt.Sprintf(" AND passport_series = $%d", argCount))
		args = append(args, filterPeople.PassportSeries)
		argCount++
	}
	if filterPeople.PassportNumber != 0 {
		q.WriteString(fmt.Sprintf(" AND passport_number = $%d", argCount))
		args = append(args, filterPeople.PassportNumber)
		argCount++
	}
	if filterPeople.Surname != "" {
		q.WriteString(fmt.Sprintf(" AND surname = $%d", argCount))
		args = append(args, filterPeople.Surname)
		argCount++
	}
	if filterPeople.Name != "" {
		q.WriteString(fmt.Sprintf(" AND name = $%d", argCount))
		args = append(args, filterPeople.Name)
		argCount++
	}
	if filterPeople.Patronymic != "" {
		q.WriteString(fmt.Sprintf(" AND patronymic = $%d", argCount))
		args = append(args, filterPeople.Patronymic)
		argCount++
	}
	if filterPeople.Address != "" {
		q.WriteString(fmt.Sprintf(" AND address = $%d", argCount))
		args = append(args, filterPeople.Address)
		argCount++
	}

	// Пагинация
	// Если limit равен 0 - OFFSET не добавляется.
	if limit > 0 {
		q.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1))
		args = append(args, limit, offset)
	} else {

		if offset > 0 {
			q.WriteString(fmt.Sprintf(" OFFSET $%d", argCount))
			args = append(args, offset)
		}
	}

	query := q.String()

	stmt, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("prepare error: %w, operation: %s", err, op)
	}

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("query error: %w, operation: %s", err, op)
	}

	var peopleList []entities.People

	for rows.Next() {
		var people entities.People
		if err := rows.Scan(&people.ID, &people.PassportSeries, &people.PassportNumber, &people.Surname, &people.Name, &people.Patronymic, &people.Address); err != nil {
			return nil, fmt.Errorf("scan error: %w, operation: %s", err, op)
		}
		peopleList = append(peopleList, people)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w, operation: %s", err, op)
	}

	return peopleList, nil
}

func (p *PeopleManagePostgres) Update(ctx context.Context, people entities.People) error {
	const op = "postgres.People.Update"

	// Проверяем, что ID предоставлен и хотя бы одно значение для обновления задано
	if people.ID == 0 {
		return fmt.Errorf("missing ID, operation: %s", op)
	}
	if people.PassportSeries == 0 && people.PassportNumber == 0 && people.Surname == "" &&
		people.Name == "" && people.Patronymic == "" && people.Address == "" {
		return fmt.Errorf("no values to update, operation: %s", op)
	}

	// Конструктор строки для запроса
	var q strings.Builder
	q.WriteString(`UPDATE people_info SET`)

	var args []interface{}
	argCount := 1

	// Добавление значений в запрос
	if people.PassportSeries != 0 {
		q.WriteString(fmt.Sprintf(" passport_series = $%d,", argCount))
		args = append(args, people.PassportSeries)
		argCount++
	}
	if people.PassportNumber != 0 {
		q.WriteString(fmt.Sprintf(" passport_number = $%d,", argCount))
		args = append(args, people.PassportNumber)
		argCount++
	}
	if people.Surname != "" {
		q.WriteString(fmt.Sprintf(" surname = $%d,", argCount))
		args = append(args, people.Surname)
		argCount++
	}
	if people.Name != "" {
		q.WriteString(fmt.Sprintf(" name = $%d,", argCount))
		args = append(args, people.Name)
		argCount++
	}
	if people.Patronymic != "" {
		q.WriteString(fmt.Sprintf(" patronymic = $%d,", argCount))
		args = append(args, people.Patronymic)
		argCount++
	}
	if people.Address != "" {
		q.WriteString(fmt.Sprintf(" address = $%d", argCount))
		args = append(args, people.Address)
		argCount++
	}

	// Удаляем последнюю запятую
	query := q.String()
	if len(query) > len("UPDATE people_info SET") {
		query = query[:len(query)-1] // Удаление последней запятой
	}

	// Добавление ID обновляемой записи
	query += fmt.Sprintf(" WHERE id = $%d", argCount)
	args = append(args, people.ID)

	stmt, err := p.db.PrepareContext(ctx, query)
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
		return fmt.Errorf("no rows updated, operation: %s", op)
	}

	return nil
}

func (p *PeopleManagePostgres) List(ctx context.Context) ([]entities.People, error) {
	const op = "postgres.People.List"

	q := `SELECT id, passport_series, passport_number, surname, name, patronymic, address FROM people_info;`

	stmt, err := p.db.PrepareContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("prepare error: %w, operation: %s", err, op)
	}

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("query error: %w, operation: %s", err, op)
	}

	var peopleList []entities.People

	for rows.Next() {
		var people entities.People
		if err := rows.Scan(&people.ID, &people.PassportSeries, &people.PassportNumber, &people.Surname, &people.Name, &people.Patronymic, &people.Address); err != nil {
			return nil, fmt.Errorf("scan error: %w, operation: %s", err, op)
		}
		peopleList = append(peopleList, people)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w, operation: %s", err, op)
	}

	return peopleList, nil
}

func (p *PeopleManagePostgres) Delete(ctx context.Context, peopleID int) error {
	const op = "postgres.People.Delete"

	// Удаляется только пользователь, остальные таблицы не затрагиваются.
	// Например, в дальнейшем, это позволит поменять исполнителя задачи.
	// Foreign key для time_entries с опцией ON DELETE SET NULL.
	q := `DELETE FROM people_info WHERE id = $1;`

	stmt, err := p.db.PrepareContext(ctx, q)
	if err != nil {
		return fmt.Errorf("prepare error: %w, operation: %s", err, op)
	}

	result, err := stmt.ExecContext(ctx, peopleID)
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
