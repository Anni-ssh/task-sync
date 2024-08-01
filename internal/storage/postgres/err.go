package postgres

import "errors"

var (
	ErrInputData      = errors.New("incorrect input data")
	ErrNoRecordsFound = errors.New("no records found")
)
