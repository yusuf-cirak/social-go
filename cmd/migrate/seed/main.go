package main

import (
	"github.com/yusuf-cirak/social/internal/db"
	"github.com/yusuf-cirak/social/internal/db/seed"
	"github.com/yusuf-cirak/social/internal/env"
	"github.com/yusuf-cirak/social/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "")
	if addr == "" {
		panic("DB_ADDR environment variable is not set")
	}

	conn, err := db.New(addr, 25, 25, "5m")
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)

	seed.Seed(store)
}
