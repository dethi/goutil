package sqlstore

import (
	"context"
	"database/sql"
	"sync"
)

// WrapWithStmtCache wraps a PrepareExecutor and returns an Executor. The
// returned executor will automatically prepare SQL queries and reuse them. The
// statement cache can be disabled on a per-execution based using the
// WithStmtCacheDisabled function.
//
// You *must* disable the cache if the SQL query is not constant, i.e. if the
// query is built using string concatenation. If you don't disable the cache,
// both the service and MySQL will leak memory.
//
// For example, if the query is using the IN keyword, there is a high probability
// that the query is built using string concatenation.
func WrapWithStmtCache(store PrepareExecutor) Executor {
	return &stmtCacheExecutor{store: store}
}

type stmtCacheExecutor struct {
	store     PrepareExecutor
	stmtcache sync.Map
}

// ExecContext executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (e *stmtCacheExecutor) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	stmt, ok := e.prepare(ctx, query)
	if !ok {
		// Fallback if we weren't able to load/create the statement
		return e.store.ExecContext(ctx, query, args...)
	}
	return stmt.ExecContext(ctx, args...)
}

// QueryContext executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (e *stmtCacheExecutor) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	stmt, ok := e.prepare(ctx, query)
	if !ok {
		// Fallback if we weren't able to load/create the statement
		return e.store.QueryContext(ctx, query, args...)
	}
	return stmt.QueryContext(ctx, args...)
}

// QueryRowContext executes a query that is expected to return at most one row.
// QueryRowContext always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards
// the rest.
func (e *stmtCacheExecutor) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	stmt, ok := e.prepare(ctx, query)
	if !ok {
		// Fallback if we weren't able to load/create the statement
		return e.store.QueryRowContext(ctx, query, args...)
	}
	return stmt.QueryRowContext(ctx, args...)
}

func (e *stmtCacheExecutor) prepare(ctx context.Context, query string) (stmt *sql.Stmt, ok bool) {
	if isStmtCacheDisabled(ctx) {
		return nil, false
	}

	stmt, ok = e.cacheLoad(query)
	if !ok {
		var err error

		stmt, err = e.store.PrepareContext(ctx, query)
		if err != nil {
			return nil, false
		}
		e.cacheStore(query, stmt)
	}
	return stmt, true
}

func (e *stmtCacheExecutor) cacheLoad(query string) (*sql.Stmt, bool) {
	var stmt *sql.Stmt

	val, ok := e.stmtcache.Load(query)
	if ok {
		stmt = val.(*sql.Stmt)
	}
	return stmt, ok
}

func (e *stmtCacheExecutor) cacheStore(query string, stmt *sql.Stmt) {
	e.stmtcache.Store(query, stmt)
}

type ctxkey int

const (
	disableStmtCache ctxkey = iota
)

// WithStmtCacheDisabled sets a flag to disable the statement cache. The cache
// will be disabled for every function call using the returned context.
func WithStmtCacheDisabled(ctx context.Context) context.Context {
	return context.WithValue(ctx, disableStmtCache, struct{}{})
}

func isStmtCacheDisabled(ctx context.Context) bool {
	return ctx.Value(disableStmtCache) != nil
}
