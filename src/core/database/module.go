package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/frostgrove/vv-template/src/config"
	"github.com/frostgrove/vv/crud"
	"github.com/frostgrove/vv/crud/adapter/crudsql"
	"github.com/frostgrove/vv/utils/vvdb"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/fx"
)

// Register opens the application's database/sql pool once and exposes a vv
// source over it. No generated repository depends on a particular ORM or
// driver; it only asks Fx for this crud.Source.
func Register() fx.Option {
	return fx.Options(
		fx.Provide(Open, Source),
		// The transport may not depend on a repository yet, but invalid database
		// configuration should still fail application start rather than the first
		// request that happens to need one.
		fx.Invoke(func(crud.Source) {}),
	)
}

func Open(cfg *config.Config, lc fx.Lifecycle) (*sql.DB, error) {
	db, err := vvdb.Open(cfg.Infra.Db)
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(context.Background()); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("connecting to database: %w", err)
	}
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error { return db.Close() },
	})
	return db, nil
}

func Source(cfg *config.Config, db *sql.DB) (crud.Source, error) {
	switch cfg.Infra.Db.Engine {
	case vvdb.Postgres:
		return crudsql.Postgres(db), nil
	case vvdb.MySQL:
		return crudsql.MySQL(db), nil
	case vvdb.MariaDB:
		return crudsql.MariaDB(db), nil
	case vvdb.SQLite:
		return crudsql.SQLite(db), nil
	default:
		return nil, fmt.Errorf("unsupported database engine %q", cfg.Infra.Db.Engine)
	}
}
