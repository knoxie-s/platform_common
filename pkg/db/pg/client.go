package pg

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/knoxie-s/platform_common/pkg/db"
)

type pgClient struct {
	masterDBC db.DB
}

// NewClient create dbClient
func NewClient(ctx context.Context, dsn string) (db.Client, error) {
	conn, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return &pgClient{
		masterDBC: New(conn),
	}, nil
}

// DB Get db contract
func (cl *pgClient) DB() db.DB {
	return cl.masterDBC
}

// Close db connection
func (cl *pgClient) Close() error {
	if cl.masterDBC != nil {
		cl.masterDBC.Close()
	}

	return nil
}
