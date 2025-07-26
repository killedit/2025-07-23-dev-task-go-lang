package main

import (
	"fmt"
	"testing"
)

func TestSchrodingerStore_BasicOperations(t *testing.T) {
	db := GetDB()
	defer db.Close()

	store := NewSchrodingerStore(db)

	err := store.Put("test_key", "test_value")
	if err != nil {
		t.Errorf("Put operation failed: %v", err)
	}

	value, err := store.Get("test_key")
	if err != nil {
		t.Errorf("Get operation failed: %v", err)
	}
	if value != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", value)
	}

	err = store.Delete("test_key")
	if err != nil {
		t.Errorf("Delete operation failed: %v", err)
	}

	_, err = store.Get("test_key")
	if err == nil {
		t.Error("Expected error after deletion, but got none")
	}
}

func TestSchrodingerStore_QuantumChaos(t *testing.T) {
	db := GetDB()
	defer db.Close()

	store := NewSchrodingerStore(db)

	testData := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	for key, value := range testData {
		store.Put(key, value)
	}

	chaosCount := 0
	totalTests := 100

	for i := 0; i < totalTests; i++ {
		value, err := store.Get("key1")
		if err != nil {
			continue
		}
		if value != "value1" {
			chaosCount++
		}
	}

	chaosPercentage := float64(chaosCount) / float64(totalTests) * 100
	t.Logf("Quantum chaos observed in %.1f%% of Get operations", chaosPercentage)

	if chaosPercentage < 10 || chaosPercentage > 50 {
		t.Logf("Warning: Quantum chaos percentage (%.1f%%) is outside expected range (10-50%%)", chaosPercentage)
	}
}

func TestSchrodingerStore_Dump(t *testing.T) {
	db := GetDB()
	defer db.Close()

	store := NewSchrodingerStore(db)
	store.Put("dump_test_key", "dump_test_value")
	store.Dump()
}

func TestSchrodingerStore_ConcurrentOperations(t *testing.T) {
	db := GetDB()
	defer db.Close()

	store := NewSchrodingerStore(db)
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			key := fmt.Sprintf("concurrent_key_%d", id)
			value := fmt.Sprintf("concurrent_value_%d", id)

			err := store.Put(key, value)
			if err != nil {
				t.Errorf("Concurrent Put failed: %v", err)
			}

			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("concurrent_key_%d", i)
		expectedValue := fmt.Sprintf("concurrent_value_%d", i)

		value, err := store.Get(key)
		if err != nil {
			t.Errorf("Failed to get %s: %v", key, err)
			continue
		}

		if value != expectedValue {
			t.Logf("Quantum chaos: expected '%s' for key '%s', got '%s'", expectedValue, key, value)
		}
	}
}

func TestSchrodingerStore_EdgeCases(t *testing.T) {
	db := GetDB()
	defer db.Close()

	store := NewSchrodingerStore(db)

	err := store.Put("", "empty_key_value")
	if err != nil {
		t.Errorf("Put with empty key failed: %v", err)
	}

	err = store.Put("empty_value_key", "")
	if err != nil {
		t.Errorf("Put with empty value failed: %v", err)
	}

	longKey := string(make([]byte, 1000))
	longValue := string(make([]byte, 1000))

	err = store.Put(longKey, longValue)
	if err != nil {
		t.Errorf("Put with long key/value failed: %v", err)
	}

	specialKey := "key_with_特殊字符_🎉_émojis"
	specialValue := "value_with_特殊字符_🎉_émojis"

	err = store.Put(specialKey, specialValue)
	if err != nil {
		t.Errorf("Put with special characters failed: %v", err)
	}
}
func BenchmarkSchrodingerStore_Put(b *testing.B) {
	db := GetDB()
	defer db.Close()

	store := NewSchrodingerStore(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("bench_key_%d", i)
		value := fmt.Sprintf("bench_value_%d", i)
		store.Put(key, value)
	}
}

func BenchmarkSchrodingerStore_Get(b *testing.B) {
	db := GetDB()
	defer db.Close()

	store := NewSchrodingerStore(db)

	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("bench_key_%d", i)
		value := fmt.Sprintf("bench_value_%d", i)
		store.Put(key, value)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("bench_key_%d", i%1000)
		store.Get(key)
	}
}
