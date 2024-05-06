package db

import (
	"errors"
	"sync"
)

var Cache *CacheMap

type CacheMap struct {
	mu  *sync.RWMutex
	Map map[string]interface{}
}

func NewCacheMap() {
	// return
	Cache = &CacheMap{
		Map: make(map[string]interface{}),
		mu:  &sync.RWMutex{},
	}
}

func (c *CacheMap) SetCustomCache(key string, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(key) == 0 || value == nil {
		return errors.New("key and value must not be empty")
	}
	c.Map[key] = value
	return nil

}

func (c *CacheMap) GetCustomCache(key string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists := c.Map[key]
	if !exists {
		return "", errors.New("error getting cache key")
	}
	return value, nil
}

func (c *CacheMap) SetCustomCacheWithMap(data map[string]interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, value := range data {
		if key == "" {
			return errors.New("key must not be empty")
		}
		c.Map[key] = value
	}

	return nil
}
