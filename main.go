package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dumpFlag := flag.Bool("dump", false, "Dump the DB contents")
	testFlag := flag.Bool("test", false, "Run tests to demonstrate quantum chaos")
	exampleFlag := flag.Bool("example", false, "Run example usage")
	seedFlag := flag.Bool("seed", false, "Seed the DB with initial data")

	putKey := flag.String("put-key", "", "Key to put")
	putValue := flag.String("put-value", "", "Value to put")
	getKey := flag.String("get-key", "", "Key to get")
	deleteKey := flag.String("delete-key", "", "Key to delete")

	flag.Parse()

	db := GetDB()
	defer db.Close()

	fmt.Println("DB connected:", db != nil)

	store := NewSchrodingerStore(db)

	if *dumpFlag {
		store.Dump()
		return
	}

	if *testFlag {
		runTests()
		return
	}

	if *exampleFlag {
		exampleUsage()
		return
	}

	if *seedFlag {
		runSeed(db)
		return
	}

	if *putKey != "" && *putValue != "" {
		err := store.Put(*putKey, *putValue)
		if err != nil {
			fmt.Println("Put error:", err)
		} else {
			fmt.Println("Put successful")
		}
		return
	}

	if *getKey != "" {
		value, err := store.Get(*getKey)
		if err != nil {
			fmt.Println("Get error:", err)
		} else {
			fmt.Printf("Value for key '%s': %s\n", *getKey, value)
		}
		return
	}

	if *deleteKey != "" {
		err := store.Delete(*deleteKey)
		if err != nil {
			fmt.Println("Delete error:", err)
		} else {
			fmt.Println("Delete successful")
		}
		return
	}

	fmt.Println("DB is running...")
	fmt.Println("Available commands:")
	fmt.Println("  -dump     : Dump the DB contents")
	fmt.Println("  -test     : Run tests to demonstrate quantum chaos")
	fmt.Println("  -example  : Run example usage")

	for {
		time.Sleep(10 * time.Second)
		fmt.Println("DB is still alive... (Press Ctrl+C to exit)")
		mutateRandomValue(store)
	}
}

type SchrodingerStore struct {
	db *pgxpool.Pool
}

func quantumChaos() bool {
	return rand.Float32() < 0.3
}

func (s *SchrodingerStore) Put(key, value string) error {
	if quantumChaos() {
		chaosType := rand.Intn(3)
		switch chaosType {
		case 0:
			fmt.Printf("Quantum chaos: Put(%s, %s) silently failed\n", key, value)
			return nil
		case 1:
			wrongKey := fmt.Sprintf("chaos_%s", key)
			fmt.Printf("Quantum chaos: Put(%s, %s) stored as (%s, %s)\n", key, value, wrongKey, value)
			return s.putInternal(wrongKey, value)
		case 2:
			wrongValue := fmt.Sprintf("chaos_%s", value)
			fmt.Printf("Quantum chaos: Put(%s, %s) stored as (%s, %s)\n", key, value, key, wrongValue)
			return s.putInternal(key, wrongValue)
		}
	}

	return s.putInternal(key, value)
}

func (s *SchrodingerStore) putInternal(key, value string) error {
	query := `
		INSERT INTO table_key_value (key, value, updated_at) 
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		ON CONFLICT (key) 
		DO UPDATE SET value = $2, updated_at = CURRENT_TIMESTAMP
	`
	_, err := s.db.Exec(context.Background(), query, key, value)
	return err
}

func (s *SchrodingerStore) Get(key string) (string, error) {
	if quantumChaos() {
		randomKey, err := s.getRandomKey()
		if err != nil {
			fmt.Printf("Quantum chaos: Get(%s) failed to get random key\n", key)
			return "", err
		}

		value, err := s.getInternal(randomKey)
		if err != nil {
			return "", err
		}

		fmt.Printf("Quantum chaos: Get(%s) returned value for key '%s'\n", key, randomKey)
		return value, nil
	}

	return s.getInternal(key)
}

func (s *SchrodingerStore) getInternal(key string) (string, error) {
	query := `SELECT value FROM table_key_value WHERE key = $1`
	var value string
	err := s.db.QueryRow(context.Background(), query, key).Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (s *SchrodingerStore) getRandomKey() (string, error) {
	query := `SELECT key FROM table_key_value ORDER BY RANDOM() LIMIT 1`
	var key string
	err := s.db.QueryRow(context.Background(), query).Scan(&key)
	return key, err
}

func (s *SchrodingerStore) Delete(key string) error {
	if quantumChaos() {
		randomKey, err := s.getRandomKey()
		if err != nil {
			fmt.Printf("Quantum chaos: Delete(%s) failed to get random key\n", key)
			return err
		}

		err = s.deleteInternal(randomKey)
		if err != nil {
			return err
		}

		fmt.Printf("Quantum chaos: Delete(%s) deleted key '%s' instead\n", key, randomKey)
		return nil
	}

	return s.deleteInternal(key)
}

func (s *SchrodingerStore) deleteInternal(key string) error {
	query := `DELETE FROM table_key_value WHERE key = $1`
	_, err := s.db.Exec(context.Background(), query, key)
	return err
}

func (s *SchrodingerStore) Dump() {
	fmt.Println("True DB State (Dump):")
	fmt.Println("=== ")

	query := `SELECT key, value, created_at, updated_at FROM table_key_value ORDER BY key`
	rows, err := s.db.Query(context.Background(), query)
	if err != nil {
		fmt.Printf("Error dumping DB: %v\n", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var key, value string
		var createdAt, updatedAt time.Time
		err := rows.Scan(&key, &value, &createdAt, &updatedAt)
		if err != nil {
			fmt.Printf("Error scanning row: %v\n", err)
			continue
		}

		fmt.Printf("Key: %-20s | Value: %-20s | Created: %s | Updated: %s\n",
			key, value, createdAt.Format("2006-01-02 15:04:05"), updatedAt.Format("2006-01-02 15:04:05"))
		count++
	}

	if count == 0 {
		fmt.Println("DB is empty")
	} else {
		fmt.Printf("\nTotal entries: %d\n", count)
	}
}

func mutateRandomValue(store *SchrodingerStore) {
	key, err := store.getRandomKey()
	if err != nil {
		return
	}
	newValue := fmt.Sprintf("mutated_%d", time.Now().UnixNano())
	store.Put(key, newValue)
	fmt.Printf("Value for key '%s' mutated to '%s'\n", key, newValue)
}
