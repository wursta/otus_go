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

func (c lruCache) Get(key Key) (interface{}, bool) {
	item, ok := c.items[key]
	if !ok {
		return nil, false
	}
	c.queue.MoveToFront(item)
	return item.Value, true
}

func (c lruCache) Set(key Key, value interface{}) bool {
	item, inCache := c.items[key]
	if inCache {
		item.Value = value
		c.queue.MoveToFront(item)
	} else {
		newItem := c.queue.PushFront(value)
		c.items[key] = newItem
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
