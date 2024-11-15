package pg

import (
	"context"
	"fmt"
	"log"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/knoxie-s/platform_common/pkg/db"
	"github.com/knoxie-s/platform_common/pkg/db/prettier"
)

type key string

const (
	// TxKey set and get [into, from] context
	TxKey key = "tx"
)

type pg struct {
	dbc *pgxpool.Pool
}

// New DB contract
func New(dbc *pgxpool.Pool) db.DB {
	return &pg{
		dbc: dbc,
	}
}

// ScanOneContext query into dest
func (p *pg) ScanOneContext(ctx context.Context, dest any, q db.Query, args ...any) error {
	logPretty(ctx, q, args...)

	row, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanOne(dest, row)
}

// ScanAllContext query into dest
func (p *pg) ScanAllContext(ctx context.Context, dest any, q db.Query, args ...any) error {
	logPretty(ctx, q, args...)

	rows, err := p.QueryContext(ctx, q, args)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dest, rows)
}

// ExecContext query
func (p *pg) ExecContext(ctx context.Context, q db.Query, args ...any) (pgconn.CommandTag, error) {
	logPretty(ctx, q, args...)

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Exec(ctx, q.QueryRaw, args...)
	}

	return p.dbc.Exec(ctx, q.QueryRaw, args...)
}

// QueryContext create query rows
func (p *pg) QueryContext(ctx context.Context, q db.Query, args ...any) (pgx.Rows, error) {
	logPretty(ctx, q, args...)

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Query(ctx, q.QueryRaw, args...)
	}

	return p.dbc.Query(ctx, q.QueryRaw, args...)
}

// QueryRowContext create query row
func (p *pg) QueryRowContext(ctx context.Context, q db.Query, args ...any) pgx.Row {
	logPretty(ctx, q, args...)

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.QueryRow(ctx, q.QueryRaw, args...)
	}

	return p.dbc.QueryRow(ctx, q.QueryRaw, args...)
}

// Ping ...
func (p *pg) Ping(ctx context.Context) error {
	return p.dbc.Ping(ctx)
}

// Close connection
func (p *pg) Close() {
	p.dbc.Close()
}

// BeginTx transaction
func (p *pg) BeginTx(ctx context.Context, txOption pgx.TxOptions) (pgx.Tx, error) {
	return p.dbc.BeginTx(ctx, txOption)
}

// MakeContextTx create context transactions
func MakeContextTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}

func logPretty(ctx context.Context, q db.Query, args ...any) {
	queryPretty := prettier.Pretty(q.QueryRaw, prettier.PlaceholderDollar, args)
	log.Println(
		ctx,
		fmt.Sprintf("sql: %s", q.Name),
		fmt.Sprintf("query: %s", queryPretty),
	)
}
