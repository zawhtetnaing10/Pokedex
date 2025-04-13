package pokecache

import (
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second

	cases := []struct {
		key   string
		value []byte
	}{
		{
			key:   "https://www.example.com",
			value: []byte("testdata"),
		},
		{
			key:   "https://www.anotherexample.com",
			value: []byte("anothertestdata"),
		},
	}

	for _, c := range cases {

		// Create a new cache and add the values
		cache := NewCache(interval)
		cache.Add(c.key, c.value)

		// Check if the value is found.
		value, found := cache.Get(c.key)
		// If not found, fail the test
		if !found {
			t.Errorf("Expected to find key")
		}

		// If the value found is not the same as expected, fail the test.
		if string(value) != string(c.value) {
			t.Errorf("The values are not equal")
		}
	}
}

func TestReapLoop(t *testing.T) {
	// Establish times
	const baseTime = 5 * time.Millisecond
	const waitingTime = baseTime + 5*time.Millisecond

	// Create a new cache and add values
	cache := NewCache(baseTime)
	key := "https://www.example.com"
	value := []byte("testdata")
	cache.Add(key, value)

	// Immediately check the added values. They must be found since they're not expired
	_, foundUnexpired := cache.Get(key)
	if !foundUnexpired {
		t.Errorf("unexpired cache not found")
		return
	}

	// Wait longer than the cache's interval
	time.Sleep(waitingTime)

	// Check the added values again. They must be gone since the interval has expired.
	_, foundExpired := cache.Get(key)
	if foundExpired {
		t.Errorf("expired cache found.")
	}
}
