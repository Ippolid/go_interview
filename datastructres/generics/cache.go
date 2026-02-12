package generics

import (
	"errors"
	"sync"
	"time"
)

// ошибки
var (
	ErrNotFound = errors.New("no content")
	ErrExpired  = errors.New("expired")
)

// интерфейс кеша
type Cache[T any] interface {
	Get(key string) (T, error)
	Set(key string, val T) error
}

// Значение хранимое в кеше
type Value[T any] struct {
	val T
	exp time.Time
}

// Конкретная имплементация кеша
type CacheImpl[T any] struct {
	storage map[string]Value[T]
	mu      sync.RWMutex
	ttl     time.Duration
}

// создание структуры
func NewCache[T any](ttl time.Duration, cleanInterval time.Duration) Cache[T] {
	cache := &CacheImpl[T]{
		storage: make(map[string]Value[T]),
		mu:      sync.RWMutex{},
		ttl:     ttl,
	}

	//очистка мусора
	go cache.startGC(cleanInterval)

	return cache
}

// устанавливаем по ключу значени
func (c *CacheImpl[T]) Set(key string, val T) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.storage[key] = Value[T]{
		val: val,
		exp: time.Now().Add(c.ttl),
	}
	return nil
}

// получаем значение и проверяем его пригодность
func (c *CacheImpl[T]) Get(key string) (T, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var zero T
	val, ok := c.storage[key]
	if !ok {
		return zero, ErrNotFound
	}

	if time.Now().After(val.exp) {
		delete(c.storage, key)
		return zero, ErrExpired
	}

	return val.val, nil
}

// очистка мусора
func (c *CacheImpl[T]) startGC(interval time.Duration) {
	// Тикер будет "звенеть" каждые interval времени
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		// Блокируем всё хранилище на запись
		c.mu.Lock()

		// Быстро пробегаем и удаляем старье
		for key, val := range c.storage {
			if time.Now().After(val.exp) {
				delete(c.storage, key)
			}
		}

		c.mu.Unlock()
	}
}
