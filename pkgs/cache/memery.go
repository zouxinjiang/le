/*
 * Copyright (c) 2019.
 */

package cache

import (
	"sync"
	"time"
)

type memeryItem struct {
	value   interface{}
	forever bool //永久存储
	expired time.Time
}

type MemoryCache struct {
	sync.RWMutex
	data map[string]memeryItem
}

func (m *MemoryCache) Set(key string, value interface{}, timeout time.Duration) error {
	var item = memeryItem{
		value:   value,
		expired: time.Now().Add(timeout),
	}
	if timeout == time.Duration(0) {
		item.forever = true
	}
	m.Lock()
	m.data[key] = item
	m.Unlock()
	return nil
}

func (m *MemoryCache) Get(key string) interface{} {
	m.GC()
	m.RLock()
	item, ok := m.data[key]
	m.RUnlock()
	if ok {
		return item.value
	}
	return nil
}

func (m *MemoryCache) Delete(key string) error {
	m.Lock()
	delete(m.data, key)
	m.Unlock()
	return nil
}

func (m *MemoryCache) IsExists(key string) bool {
	m.GC()
	m.RLock()
	_, ok := m.data[key]
	m.RUnlock()
	return ok
}

func (m *MemoryCache) ClearAll() error {
	m.Lock()
	m.data = map[string]memeryItem{}
	m.Unlock()
	return nil
}

func (m *MemoryCache) SetOption(option CacheOption) {
}

func (m *MemoryCache) GC() {
	m.Lock()
	for k, v := range m.data {
		if !v.forever {
			if v.expired.Unix() <= time.Now().Unix() {
				//过期回收
				delete(m.data, k)
			}
		}
	}
	m.Unlock()
}

func NewMemoryCache() CacheAble {
	return &MemoryCache{
		data: map[string]memeryItem{},
	}
}

func init() {
	Register("memory", NewMemoryCache)
}
