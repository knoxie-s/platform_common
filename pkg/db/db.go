package db

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// TxFnHandler to handle repo handlers
type TxFnHandler func(ctx context.Context) error

// Client wrapper upon pgx
type Client interface {
	DB() DB
	Close() error
}

// Query struct
type Query struct {
	Name     string
	QueryRaw string
}

// SQLExecer ...
type SQLExecer interface {
	NamedExecer
	QueryExecer
}

// NamedExecer for named queries
type NamedExecer interface {
	ScanOneContext(ctx context.Context, dest any, query Query, args ...any) error
	ScanAllContext(ctx context.Context, dest any, query Query, args ...any) error
}

// QueryExecer for exec queries
type QueryExecer interface {
	ExecContext(ctx context.Context, query Query, args ...any) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, query Query, args ...any) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, query Query, args ...any) pgx.Row
}

// TxManager manager to handle transactions
type TxManager interface {
	ReadCommitted(ctx context.Context, fn TxFnHandler) error
}

// Transactor for transactions
type Transactor interface {
	BeginTx(ctx context.Context, options pgx.TxOptions) (pgx.Tx, error)
}

// Pinger connection
type Pinger interface {
	Ping(ctx context.Context) error
}

// DB contract for db actions
type DB interface {
	SQLExecer
	Pinger
	Transactor
	Close()
}
