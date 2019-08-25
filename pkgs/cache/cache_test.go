package cache

import (
	"testing"
	"time"
)

func TestMemoryCache(t *testing.T) {
	datacase := []struct {
		Key         string
		Value       interface{}
		Expire      time.Duration
		ExceptValue interface{}
	}{
		{"a1", "a1value", 0, "a1value"},
		{"a2", 2, 0, 2},
		{"a3", "a3value", time.Second * 5, "a3value"},
		{"a4", nil, 0, nil},
	}

	cache := NewCache("file")
	for _, item := range datacase {
		cache.Set(item.Key, item.Value, item.Expire)
	}
	for _, item := range datacase {
		if val := cache.Get(item.Key); val != item.ExceptValue {
			t.Log("key:", item.Key, " except value ", item.ExceptValue, " but ", val)
		}
	}
	for _, item := range datacase {
		val := cache.Get(item.Key)
		exist := cache.IsExists(item.Key)
		t.Log("key:", item.Key, "exists:", exist, " value:", val)
	}
	cache.Delete("a4")
	time.Sleep(time.Second * 5)
	for _, item := range datacase {
		val := cache.Get(item.Key)
		exist := cache.IsExists(item.Key)
		t.Log("key:", item.Key, "exists:", exist, " value:", val)
	}
}
