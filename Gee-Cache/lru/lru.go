package lru

import "container/list"

// Cache is not concurrently safe
type LRUCache struct {
	maxBytes int64
	nbytes   int64
	cache    map[string]*list.Element
	ll       *list.List
	// opetional and executed when an entry is purged
	OnEvicted func(key string, value Value)
}

//

type entry struct {
	key   string
	value Value
}

// count how many bytes it takes
type Value interface {
	Len() int
}

// LRUCache constructor
func New(maxBytes int64, OnEvicted func(string, Value)) *LRUCache {
	return &LRUCache{
		maxBytes:  maxBytes,
		cache:     make(map[string]*list.Element),
		ll:        list.New(),
		OnEvicted: OnEvicted,
	}
}

// add a value to the cache
func (c *LRUCache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		// if value already exist, update it
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

func (c *LRUCache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true

	}
	return nil, false
}

func (c *LRUCache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Len the number of cache entries
func (c *LRUCache) Len() int {
	return c.ll.Len()
}
