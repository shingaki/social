package main

import (
	"SOCIAL/internal/db"
	"SOCIAL/internal/env"
	"SOCIAL/internal/store"
	"log"
)

type sdConfig struct {
	conn seedConfig
}

type seedConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type config struct {
	conn dbConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func main() {
	cfg := config{
		conn: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}
	conn, err := db.New(
		cfg.conn.addr,
		cfg.conn.maxOpenConns,
		cfg.conn.maxIdleConns,
		cfg.conn.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}

	defer conn.Close()
	log.Println("database connection pool established")
	log.Println()

	store := store.NewStorage(conn)
	db.Seed(store)
}
