package main

import (
	"github.com/yusuf-cirak/social/internal/db"
	"github.com/yusuf-cirak/social/internal/env"
	"github.com/yusuf-cirak/social/internal/store"
)

const version = "1.0.0"

func main() {
	cfg := config{addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "5m"),
		},
		env: env.GetString("ENV", "development"),
	}

	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	store := store.NewStorage(db)

	app := application{config: cfg, store: store}

	mux := app.mount()
	if err := app.run(mux); err != nil {
		panic(err)
	}
}
