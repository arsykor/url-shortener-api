package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
	"url-shortener-api/internal/config"
)

const (
	errConnectionFailed = "connection to the database is failed"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewClient(ctx context.Context, attempts int, config *config.Config) (*pgxpool.Pool, error) {
	cs := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", config.Storage.Username, config.Storage.Password, config.Storage.Host, config.Storage.Port, config.Storage.Database)

	for attempts > 0 {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		pool, err := pgxpool.Connect(ctx, cs)
		if err != nil {
			time.Sleep(10 * time.Second)
			attempts--
		} else {
			return pool, nil
		}
	}
	return nil, errors.New(errConnectionFailed)
}
