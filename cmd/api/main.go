package main

import (
	"time"

	"github.com/yusuf-cirak/social/internal/auth"
	"github.com/yusuf-cirak/social/internal/db"
	"github.com/yusuf-cirak/social/internal/env"
	"github.com/yusuf-cirak/social/internal/store"
	"go.uber.org/zap"
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
		auth: authConfig{
			Secret:         env.GetString("JWT_SECRET", "dev-secret-change"),
			Issuer:         env.GetString("JWT_ISSUER", "social-go"),
			Audience:       env.GetString("JWT_AUDIENCE", "social-users"),
			AccessTokenTTL: time.Duration(env.GetInt("JWT_TTL_MINUTES", 60)) * time.Minute,
		},
	}

	//Logger

	logger := zap.Must(zap.NewProduction()).Sugar()

	defer logger.Sync() // flushes buffer, if any
	// Db

	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatalw("Failed to connect to database", "error", err)
	}

	defer db.Close()

	store := store.NewStorage(db)

	jwtMgr := auth.NewManager(cfg.auth.Secret, cfg.auth.Issuer, cfg.auth.Audience)

	app := application{config: cfg, store: store, logger: logger, jwt: jwtMgr}

	mux := app.mount()
	if err := app.run(mux); err != nil {
		logger.Fatalw("Failed to start server", "error", err)
	}
}
