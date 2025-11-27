package persistence

import (
	"context"
	"database/sql"
)

type Executer interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) (*sql.Row, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}
