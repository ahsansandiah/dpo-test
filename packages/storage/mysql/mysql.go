package mysqlDatabase

import (
	"database/sql"
	"time"

	"github.com/ahsansandiah/dpo-test/packages/config"
	"github.com/cenkalti/backoff"
	_ "github.com/go-sql-driver/mysql"
)

type MySQL interface {
	Connect() (*sql.DB, error)
}

type Options struct {
	driver  string
	dns     string
	maxOpen int
	maxIdle int
}

func NewMySQL(cfg *config.Config) *Options {
	opt := new(Options)
	opt.driver = cfg.DatabaseDriver
	opt.dns = cfg.DatabaseDNS
	opt.maxOpen = cfg.DatabaseMaxOpenConnections
	opt.maxIdle = cfg.DatabaseMaxIdleConnections

	return opt
}

func (o *Options) Connect() (*sql.DB, error) {
	database, err := sql.Open(o.driver, o.dns)
	if err != nil {
		return nil, err
	}

	database.SetMaxOpenConns(o.maxOpen)
	database.SetMaxIdleConns(o.maxIdle)
	database.SetConnMaxLifetime(time.Hour)

	if err := backoff.Retry(func() error {
		if err := database.Ping(); err != nil {
			return err
		}

		return nil
	}, backoff.NewExponentialBackOff()); err != nil {
		return nil, err
	}

	return database, nil
}
