package sqlstoretest

import (
	"context"
	"database/sql"
)

type MockExecutor struct {
	ExecFn     func(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryFn    func(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowFn func(ctx context.Context, query string, args ...interface{}) *sql.Row

	PrepareFn func(ctx context.Context, query string) (*sql.Stmt, error)
}

func (e *MockExecutor) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return e.ExecFn(ctx, query, args...)
}

func (e *MockExecutor) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return e.QueryFn(ctx, query, args...)
}

func (e *MockExecutor) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return e.QueryRowFn(ctx, query, args...)
}

func (e *MockExecutor) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return e.PrepareFn(ctx, query)
}
