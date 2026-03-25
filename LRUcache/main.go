package main

import (
	"sync"
	"time"
)

// think of it as Node
type entry struct {
	key       string
	value     interface{}
	expiresAt time.Time
	prev      *entry
	next      *entry
}

type LRUCache struct {
	capacity int
	items    map[string]*entry
	head     *entry
	tail     *entry
	mu       sync.Mutex
}

func NewLRUCache(cap int) *LRUCache {
	return &LRUCache{
		capacity: cap,
		items:    make(map[string]*entry),
	}
}

func (c *LRUCache) moveToFront(e *entry) {
	if c.head == e {
		return
	}

	c.remove(e)
	c.addToFront(e)
}

func (c *LRUCache) addToFront(e *entry) {
	e.next = c.head
	e.prev = nil

	if c.head != nil {
		c.head.prev = e
	}
	c.head = e

	if c.tail == nil {
		c.tail = e
	}
}

func (c *LRUCache) remove(e *entry) {
	if e.prev != nil {
		e.prev.next = e.next
	} else {
		c.head = e.next
	}

	if e.next != nil {
		e.next.prev = e.prev
	} else {
		c.tail = e.prev
	}
}

func (e *entry) isExpired() bool {
	return !e.expiresAt.IsZero() && time.Now().After(e.expiresAt)
}

func main() {

}
