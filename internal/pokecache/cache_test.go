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
		cache := NewCache(interval)

		cache.Add(c.key, c.value)

		value, found := cache.Get(c.key)
		if !found {
			t.Errorf("Expected to find key")
		}

		if string(value) != string(c.value) {
			t.Errorf("The values are not equal")
		}
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitingTime = baseTime + 5*time.Millisecond

	cache := NewCache(baseTime)
	key := "https://www.example.com"
	value := []byte("testdata")
	cache.Add(key, value)

	_, foundUnexpired := cache.Get(key)
	if !foundUnexpired {
		t.Errorf("unexpired cache not found")
		return
	}

	time.Sleep(waitingTime)

	_, foundExpired := cache.Get(key)
	if foundExpired {
		t.Errorf("expired cache found.")
	}
}
