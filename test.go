package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func runTests() {
	fmt.Println("===")
	fmt.Println("Testing the DB")
	fmt.Println("===")

	db := GetDB()
	defer db.Close()

	store := NewSchrodingerStore(db)
	rand.Seed(time.Now().UnixNano())

	fmt.Println("\n#1. Testing Put operations:")
	testData := map[string]string{
		"cat":  "meow",
		"dog":  "woof",
		"bird": "tweet",
		"fish": "blub",
		"cow":  "moo",
	}

	for key, value := range testData {
		fmt.Printf("Putting %s = %s\n", key, value)
		err := store.Put(key, value)
		if err != nil {
			log.Printf("Error putting %s: %v", key, err)
		}
	}

	fmt.Println("\n#2. Testing Get operations:")
	for key := range testData {
		value, err := store.Get(key)
		if err != nil {
			fmt.Printf("Get(%s): Error - %v\n", key, err)
		} else {
			fmt.Printf("Get(%s): %s\n", key, value)
		}
	}

	fmt.Println("\n#3. Testing Delete operations:")
	keysToDelete := []string{"cat", "fish"}
	for _, key := range keysToDelete {
		fmt.Printf("Deleting %s\n", key)
		err := store.Delete(key)
		if err != nil {
			log.Printf("Error deleting %s: %v", key, err)
		}
	}

	fmt.Println("\n#4. Final database state:")
	store.Dump()

	fmt.Println("\n5. Trying to get deleted keys:")
	for _, key := range keysToDelete {
		value, err := store.Get(key)
		if err != nil {
			fmt.Printf("Get(%s): Error - %v\n", key, err)
		} else {
			fmt.Printf("Get(%s): %s (quantum chaos might have returned wrong data!)\n", key, value)
		}
	}
}

func exampleUsage() {
	fmt.Println("Schrödinger's Database Example Usage")
	fmt.Println("===")

	db := GetDB()
	defer db.Close()

	store := NewSchrodingerStore(db)

	fmt.Println("\nExample 1: Basic operations")
	store.Put("name", "Alice")
	store.Put("age", "30")
	store.Put("city", "New York")

	name, _ := store.Get("name")
	age, _ := store.Get("age")
	city, _ := store.Get("city")

	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Age: %s\n", age)
	fmt.Printf("City: %s\n", city)

	fmt.Println("\nExample 2: Delete operation")
	store.Delete("age")

	deletedAge, err := store.Get("age")
	if err != nil {
		fmt.Printf("Age was deleted (or quantum chaos intervened): %v\n", err)
	} else {
		fmt.Printf("Age still exists (quantum chaos!): %s\n", deletedAge)
	}

	fmt.Println("\nExample 3: True database state")
	store.Dump()
}
