package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/snowmerak/api-key-server-template/pkg/client/postgres/queries"
)

type Option struct {
	host     string
	port     int
	username string
	password string
	database string
}

func NewOption(host string, port int, username, password, database string) *Option {
	return &Option{
		host:     host,
		port:     port,
		username: username,
		password: password,
		database: database,
	}
}

func (o *Option) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", o.username, o.password, o.host, o.port, o.database)
}

type DB struct {
	conn    *pgx.Conn
	queries *queries.Queries
}

func NewDB(ctx context.Context, option *Option) (*DB, error) {
	conn, err := pgx.Connect(ctx, option.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgresql instance: %v", err)
	}

	return &DB{
		conn:    conn,
		queries: queries.New(conn),
	}, nil
}
