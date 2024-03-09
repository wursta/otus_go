package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type lruCacheEl struct {
	key   Key
	value interface{}
}

func (c lruCache) Get(key Key) (interface{}, bool) {
	item, ok := c.items[key]
	if !ok {
		return nil, false
	}
	c.queue.MoveToFront(item)
	cacheEl := item.Value.(*lruCacheEl)
	return cacheEl.value, true
}

func (c lruCache) Set(key Key, value interface{}) bool {
	item, inCache := c.items[key]

	if inCache {
		cacheEl := item.Value.(*lruCacheEl)
		cacheEl.value = value
		c.queue.MoveToFront(item)
	} else {
		if len(c.items) == c.capacity {
			backNode := c.queue.Back()
			cacheEl := backNode.Value.(*lruCacheEl)
			delete(c.items, cacheEl.key)
			c.queue.Remove(backNode)
		}

		cacheEl := &lruCacheEl{key: key, value: value}
		c.items[key] = c.queue.PushFront(cacheEl)
	}

	return inCache
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
