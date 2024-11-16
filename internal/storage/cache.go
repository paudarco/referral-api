package storage

import (
	"sync"
	"time"
)

// Cache представляет in-memory кэш на основе sync.Map
type CacheStorage struct {
	store sync.Map
}

type cacheItem struct {
	value     interface{}
	expiresAt time.Time
}

// NewCache создает новый экземпляр кэша
func NewCache() *CacheStorage {
	cache := &CacheStorage{}
	go cache.StartCleanup()
	return cache
}

// Set добавляет значение в кэш с указанным временем жизни
func (c *CacheStorage) Set(key int, value interface{}, ttl time.Time) {
	c.store.Store(key, cacheItem{
		value:     value,
		expiresAt: ttl,
	})
}

// Get получает значение из кэша
func (c *CacheStorage) Get(key int) (interface{}, bool) {
	value, ok := c.store.Load(key)
	if !ok {
		return nil, false
	}

	item := value.(cacheItem)
	if time.Now().After(item.expiresAt) {
		c.store.Delete(key)
		return nil, false
	}

	return item.value, true
}

// startCleanup запускает периодическую очистку устаревших записей
func (c *CacheStorage) StartCleanup() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		c.Cleanup()
	}
}

func (c *CacheStorage) Cleanup() {
	c.store.Range(func(key, value interface{}) bool {
		item := value.(cacheItem)
		if time.Now().After(item.expiresAt) {
			c.store.Delete(key)
		}
		return true
	})
}

func (c *CacheStorage) DeleteUserCode(id int) {
	c.store.Delete(id)
}
