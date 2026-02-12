package generics

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestCacheDifferentVal(t *testing.T) {
	tests := []struct {
		name string
		key  string
		val  any
	}{
		{
			name: "string value",
			key:  "test",
			val:  "yes",
		},
		{
			name: "int value",
			key:  "test",
			val:  11,
		},
		{
			name: "error value",
			key:  "test",
			val:  fmt.Errorf("test"),
		},
	}

	ttl := time.Second
	cache := NewCache[any](ttl, 10*ttl)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cache.Set(tt.key, tt.val)
			if err != nil {
				t.Error(err)
			}

			got, err := cache.Get(tt.key)
			if err != nil {
				t.Error(err)
			}

			if got != tt.val {
				t.Fatalf("Got %d, want %d", got, tt.val)
			}
		})
	}

}

func TestCacheExpVal(t *testing.T) {
	key := "test"
	val := "yes"
	ttl := time.Second
	cache := NewCache[any](ttl, 10*ttl)

	err := cache.Set(key, val)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Second * 2)

	got, err := cache.Get(key)
	if errors.Is(err, ErrExpired) {
		t.Log(got, err)
	}

}
