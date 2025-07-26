package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewSchrodingerStore(db *pgxpool.Pool) *SchrodingerStore {
	store := &SchrodingerStore{db: db}
	store.initTable()
	return store
}

func (s *SchrodingerStore) initTable() {
	query := `
		CREATE TABLE IF NOT EXISTS table_key_value (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := s.db.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}

func runSeed(db *pgxpool.Pool) {
	fmt.Println("Seeding the DB")
	fmt.Println("===")

	store := NewSchrodingerStore(db)

	store.Put("cat", "meow")
	store.Put("dog", "dark")
	store.Put("bird", "tweet")
	fmt.Println("Seeding complete.")
}
