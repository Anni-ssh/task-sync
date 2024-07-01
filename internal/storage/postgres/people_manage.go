package postgres

import (
	"TaskSync/internal/entities"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type PeopleManagePostgres struct {
	db *sql.DB
}

func NewPeopleManage(db *sql.DB) *PeopleManagePostgres {
	return &PeopleManagePostgres{db: db}
}

func (p *PeopleManagePostgres) Create(ctx context.Context, people entities.People) (int, error) {

	const op = "postgres.People.Create"

	q := `INSERT INTO people_info (passport_series, passport_number, surname, name, patronymic, address) 
      VALUES($1, $2, $3, $4, $5, $6) 
      RETURNING id;`

	var id int

	err := p.db.QueryRowContext(ctx, q, people.PassportSeries, people.PassportNumber, people.Surname, people.Name, people.Patronymic, people.Address).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("database error: %w, operation: %s", err, op)
	}

	return id, nil

}

func (p *PeopleManagePostgres) GetByID(ctx context.Context, peopleID int) (entities.People, error) {
	const op = "postgres.People.Get"

	q := `SELECT id, passport_series, passport_number, surname, name, patronymic, address FROM people_info WHERE id = $1;`

	var people entities.People

	rows, err := p.db.QueryContext(ctx, q, peopleID)
	if err != nil {
		return people, fmt.Errorf("database error: %w, operation: %s", err, op)
	}

	if rows.Next() {
		err := rows.Scan(&people.ID, &people.PassportSeries, &people.PassportNumber, &people.Surname, &people.Name, &people.Patronymic, &people.Address)
		if err != nil {
			return people, fmt.Errorf("scan error: %w, operation: %s", err, op)
		}
	}

	if err := rows.Err(); err != nil {
		return people, fmt.Errorf("rows error: %w, operation: %s", err, op)
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
	q.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1))
	args = append(args, limit, offset)

	query := q.String()

	var peopleList []entities.People

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("database error: %w, operation: %s", err, op)
	}

	for rows.Next() {
		var people entities.People
		err := rows.Scan(&people.ID, &people.PassportSeries, &people.PassportNumber, &people.Surname, &people.Name, &people.Patronymic, &people.Address)
		if err != nil {
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

	// Проверяем, что все значения в структуре не пустые
	if people.PassportSeries == 0 && people.PassportNumber == 0 && people.Surname == "" &&
		people.Name == "" && people.Patronymic == "" && people.Address == "" && people.ID == 0 {
		return fmt.Errorf("incorrect values or their absence, operation: %s", op)
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

	// Доабавление ID обновляемой записи
	q.WriteString(fmt.Sprintf(" WHERE id = $%d", argCount))
	args = append(args, people.ID)

	result, err := p.db.ExecContext(ctx, q.String(), args...)
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

func (p *PeopleManagePostgres) List(ctx context.Context) ([]entities.People, error) {
	const op = "postgres.People.List"

	q := `SELECT id, passport_series, passport_number, surname, name, patronymic, address FROM people_info;`

	var peopleList []entities.People

	rows, err := p.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("database error: %w, operation: %s", err, op)
	}

	for rows.Next() {
		var people entities.People
		err := rows.Scan(&people.ID, &people.PassportSeries, &people.PassportNumber, &people.Surname, &people.Name, &people.Patronymic, &people.Address)
		if err != nil {
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

	result, err := p.db.ExecContext(ctx, q, peopleID)
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
