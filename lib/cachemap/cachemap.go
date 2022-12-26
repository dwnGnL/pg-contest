package cachemap

import (
	"sync"
)

type CacheMaper[K comparable, V any] struct {
	sMap *sync.Map
}

func (c *CacheMaper[K, V]) Load(key K) (V, bool) {
	value, ok := c.sMap.Load(key)
	if !ok {
		var emtpyV V
		return emtpyV, false
	}
	return value.(V), ok
}

func (c *CacheMaper[K, V]) Store(key K, value V) {
	c.sMap.Store(key, value)
}

func (c *CacheMaper[K, V]) Range(f func(key K, value V) bool) {
	c.sMap.Range(func(keyAny, valueAny any) bool {
		return f(keyAny.(K), valueAny.(V))
	})
}

func (c *CacheMaper[K, V]) IsZero() bool {
	count := 0
	c.sMap.Range(func(_, _ any) bool {
		count++
		return true
	})
	return count == 0
}

func (c *CacheMaper[K, V]) Delete(key K) {
	c.sMap.Delete(key)
}

func (c *CacheMaper[K, V]) EquateTo(m map[K]V) {
	c.sMap.Range(func(key, value interface{}) bool {
		m[key.(K)] = value.(V)
		return true
	})
}

func (c *CacheMaper[K, V]) Reset() {
	c.sMap.Range(func(key, _ interface{}) bool {
		c.sMap.Delete(key)
		return true
	})
}

func (c *CacheMaper[K, V]) UpdateByMap(m map[K]V) {
	c.sMap.Range(func(key, _ interface{}) bool {
		if value, ok := m[key.(K)]; ok {
			c.sMap.Store(key, value)
			delete(m, key.(K))
			return true
		}
		c.sMap.Delete(key)
		return true
	})
	for key, value := range m {
		c.sMap.Store(key, value)
	}
}

func NewCacheMap[K comparable, V any]() *CacheMaper[K, V] {
	return &CacheMaper[K, V]{sMap: &sync.Map{}}
}
