package memorycache

import (
	"errors"
	"sync"

	"time"
)

// Cache struct cache
type Cache struct {
	sync.RWMutex
	items             map[string]Item
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
}

// Item struct cache item
type Item struct {
	Value      interface{}
	Expiration int64
	Created    time.Time
}

// New. Initializing a new memory cache
func New(defaultExpiration, cleanupInterval time.Duration) *Cache {

	items := make(map[string]Item)

	// cache item
	cache := Cache{
		items:             items,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}

	if cleanupInterval > 0 {
		cache.StartGC()
	}

	return &cache
}

// Set setting a cache by key
func (c *Cache) Set(key string, value interface{}, duration time.Duration) {

	var expiration int64

	if duration == 0 {
		duration = c.defaultExpiration
	}

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.Lock()

	defer c.Unlock()

	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}

}

// Get getting a cache by key
func (c *Cache) Get(key string) (interface{}, bool) {

	c.RLock()

	defer c.RUnlock()

	item, found := c.items[key]

	// cache not found
	if !found {
		return nil, false
	}

	if item.Expiration > 0 {

		// cache expired
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}

	}

	return item.Value, true
}

// Delete cache by key
// Return false if key not found
func (c *Cache) Delete(key string) error {

	c.Lock()

	defer c.Unlock()

	if _, found := c.items[key]; !found {
		return errors.New("Key not found")
	}

	delete(c.items, key)

	return nil
}

// StartGC start Garbage Collection
func (c *Cache) StartGC() {
	go c.GC()
}

// GC Garbage Collection
func (c *Cache) GC() {

	for {

		<-time.After(c.cleanupInterval)

		if c.items == nil {
			return
		}

		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearItems(keys)

		}

	}

}

// expiredKeys returns key list which are expired.
func (c *Cache) expiredKeys() (keys []string) {

	c.RLock()

	defer c.RUnlock()

	for k, i := range c.items {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}

	return
}

// clearItems removes all the items which key in keys.
func (c *Cache) clearItems(keys []string) {

	c.Lock()

	defer c.Unlock()

	for _, k := range keys {
		delete(c.items, k)
	}
}
