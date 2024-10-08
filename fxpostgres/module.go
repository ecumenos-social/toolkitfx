package fxpostgres

import (
	"context"

	postgres "github.com/ecumenos-social/postgresql-driver"
	"github.com/jackc/pgx/v4"
	"go.uber.org/fx"
)

type Config struct {
	URL            string `split_words:"true"`
	MigrationsPath string `split_words:"true"`
}

var Module = fx.Options(
	fx.Provide(func(lc fx.Lifecycle, cfg *Config) (Driver, error) {
		driver, err := postgres.New(context.Background(), cfg.URL)
		if err != nil {
			return nil, err
		}
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				return driver.Ping(ctx)
			},
			OnStop: func(context.Context) error {
				driver.Close()
				return nil
			},
		})

		return driver, nil
	}),
	fx.Provide(NewMigrationsRunner),
)

type Driver interface {
	Ping(ctx context.Context) error
	Close()
	CountRows(ctx context.Context, query string, args ...interface{}) (int, error)
	ExecuteQuery(ctx context.Context, query string, args ...interface{}) error
	QueryRow(ctx context.Context, query string, args ...interface{}) (pgx.Row, error)
	QueryRows(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
}
