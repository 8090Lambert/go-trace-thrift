package service

import "sync"

type CoSafeStorage struct {
	mu        *sync.RWMutex
	coSafeMap map[interface{}]interface{}
}

func NewCoSafeStorage() *CoSafeStorage {
	safeStorage := new(CoSafeStorage)
	safeStorage.mu = new(sync.RWMutex)
	safeStorage.coSafeMap = make(map[interface{}]interface{})
	return safeStorage
}

func (c *CoSafeStorage) Set(key, val interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.coSafeMap[key] = val
}

func (c *CoSafeStorage) Get(key interface{}) (value interface{}, result bool) {
	c.mu.RLock()
	if _, ok := c.coSafeMap[key]; !ok {
		c.mu.RUnlock()
		return nil, false
	}
	c.mu.RUnlock()
	return c.coSafeMap[key], true
}

func (c *CoSafeStorage) Delete(key interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.coSafeMap, key)
}
