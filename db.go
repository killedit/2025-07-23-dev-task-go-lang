package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var (
	dbPool *pgxpool.Pool
	once   sync.Once
)

func GetDB() *pgxpool.Pool {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Printf("No .env file loaded or error loading: %v", err)
		}

		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")

		if host == "" || port == "" || user == "" || password == "" || dbname == "" {
			log.Fatal("Database environment variables (DB_HOST, DB_PORT, etc.) must be set")
		}

		dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, password, host, port, dbname)

		pool, err := pgxpool.New(context.Background(), dsn)
		if err != nil {
			log.Fatalf("Unable to connect to database: %v", err)
		}

		if err := pool.Ping(context.Background()); err != nil {
			log.Fatalf("Unable to ping database: %v", err)
		}

		dbPool = pool
		log.Println("Connected to PostgreSQL database")
	})

	return dbPool
}
