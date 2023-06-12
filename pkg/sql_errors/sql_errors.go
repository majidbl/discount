package sql_errors

import (
	"strings"
)

type SqlNotFoundError struct{}

var SqlNotFound *SqlNotFoundError

func (s *SqlNotFoundError) Error() string {
	return "row not found"
}

func NewSqlNotFoundError() error {
	return &SqlNotFoundError{}
}

type SqlError struct{}

func NewSqlError() error {
	return &SqlError{}
}

func (s *SqlError) Error() string {
	return "sql error"
}

func ParseSqlErrors(err error) error {
	if strings.Contains(err.Error(), "no rows in result set") {
		return NewSqlNotFoundError()
	}
	return NewSqlError()
}
