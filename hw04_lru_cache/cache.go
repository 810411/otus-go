package hw04_lru_cache //nolint:golint,stylecheck
import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*Item
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(item)
		item.Value = cacheItem{key: key, value: value}
		return true
	}

	i := cacheItem{key: key, value: value}
	c.queue.PushFront(i)
	c.items[key] = c.queue.Front()

	if len(c.items) > c.capacity {
		out := c.queue.Back()
		c.queue.Remove(out)
		delete(c.items, out.Value.(cacheItem).key)
	}

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if !ok {
		return nil, false
	}

	c.queue.MoveToFront(item)

	return item.Value.(cacheItem).value, true
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	item := c.queue.Back()
	for {
		delete(c.items, item.Value.(cacheItem).key)
		item = item.Next
		if item == nil {
			break
		}
	}
	c.queue = NewList()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*Item),
	}
}
