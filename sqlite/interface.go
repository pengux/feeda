package sqlite

import "database/sql"

type (
	cruderExecer interface {
		Exec(string, ...interface{}) (sql.Result, error)
	}
	cruderQueryer interface {
		Query(string, ...interface{}) (*sql.Rows, error)
	}
	cruderQueryRower interface {
		QueryRow(string, ...interface{}) *sql.Row
	}
)
