package sqlstore_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/dethi/goutil/sqlstore"
	"github.com/dethi/goutil/sqlstore/sqlstoretest"
	"github.com/stretchr/testify/assert"
)

func TestStmtCacheExecutor(t *testing.T) {
	type args struct {
		ctx       context.Context
		query     string
		wantPanic bool // panic = a function on a prepared statement was called
	}

	ctx := context.Background()
	tests := []struct {
		name string
		args []args

		wantPrepareCount int
	}{
		{
			name: "queries are cached and reused",
			args: []args{
				args{ctx, "queryA", true},
				args{ctx, "queryB", true},
				args{ctx, "queryC", true},
				args{ctx, "queryB", true},
				args{ctx, "queryC", true},
				args{ctx, "queryA", true},
				args{ctx, "queryC", true},
				args{ctx, "queryC", true},
			},
			wantPrepareCount: 3,
		},
		{
			name: "statement cache can be disabled",
			args: []args{
				args{sqlstore.WithStmtCacheDisabled(ctx), "queryA", false},
				args{ctx, "queryB", true},
				args{ctx, "queryC", true},
				args{sqlstore.WithStmtCacheDisabled(ctx), "queryC", false},
				args{ctx, "queryB", true},
			},
			wantPrepareCount: 2,
		},
	}

	prepareCount := 0
	mockStore := &sqlstoretest.MockExecutor{
		PrepareFn: func(ctx context.Context, query string) (*sql.Stmt, error) {
			prepareCount++
			return nil, nil
		},
		ExecFn: func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
			return nil, nil
		},
		QueryFn: func(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
			return nil, nil
		},
		QueryRowFn: func(ctx context.Context, query string, args ...interface{}) *sql.Row {
			return nil
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+"/ExecContext", func(t *testing.T) {
			prepareCount = 0
			store := sqlstore.WrapWithStmtCache(mockStore)

			for _, args := range tt.args {
				assertPanics(args.wantPanic, t,
					func() { _, _ = store.ExecContext(args.ctx, args.query) },
					"ExecContext(%v)", args.query,
				)
			}

			if prepareCount != tt.wantPrepareCount {
				t.Errorf("PrepareContext() called %v, want %v", prepareCount, tt.wantPrepareCount)
			}
		})

		t.Run(tt.name+"/QueryContext", func(t *testing.T) {
			prepareCount = 0
			store := sqlstore.WrapWithStmtCache(mockStore)

			for _, args := range tt.args {
				assertPanics(args.wantPanic, t,
					func() { _, _ = store.QueryContext(args.ctx, args.query) },
					"QueryContext(%v)", args.query,
				)
			}

			if prepareCount != tt.wantPrepareCount {
				t.Errorf("PrepareContext() called %v, want %v", prepareCount, tt.wantPrepareCount)
			}
		})

		t.Run(tt.name+"/QueryRowContext", func(t *testing.T) {
			prepareCount = 0
			store := sqlstore.WrapWithStmtCache(mockStore)

			for _, args := range tt.args {
				assertPanics(args.wantPanic, t,
					func() { store.QueryRowContext(args.ctx, args.query) },
					"QueryRowContext(%v)", args.query,
				)
			}

			if prepareCount != tt.wantPrepareCount {
				t.Errorf("PrepareContext() called %v, want %v", prepareCount, tt.wantPrepareCount)
			}
		})
	}
}

func assertPanics(wantPanic bool, t assert.TestingT, f assert.PanicTestFunc, msgAndArgs ...interface{}) bool {
	if wantPanic {
		return assert.Panics(t, f, msgAndArgs...)
	}
	return assert.NotPanics(t, f, msgAndArgs...)
}

func TestStmtCacheExecutorPrepareFail(t *testing.T) {
	prepareCount := 0
	wantPrepareCount := 1

	queryRowCount := 0
	wantQueryRowCount := 1

	mockStore := &sqlstoretest.MockExecutor{
		PrepareFn: func(ctx context.Context, query string) (*sql.Stmt, error) {
			prepareCount++
			return nil, errors.New("something wrong happened")
		},
		QueryRowFn: func(ctx context.Context, query string, args ...interface{}) *sql.Row {
			queryRowCount++
			return nil
		},
	}
	store := sqlstore.WrapWithStmtCache(mockStore)
	store.QueryRowContext(context.Background(), "queryA")

	if prepareCount != wantPrepareCount {
		t.Errorf("PrepareContext() called %v, want %v", prepareCount, wantPrepareCount)
	}
	if queryRowCount != wantQueryRowCount {
		t.Errorf("PrepareContext() called %v, want %v", queryRowCount, wantQueryRowCount)
	}
}
