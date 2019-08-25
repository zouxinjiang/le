/*
 * Copyright (c) 2019.
 */

package cache

import (
	"fmt"
	"time"
)

type CacheOption interface {
	GetOption(key string) string
	SetOption(key string) string
}

type CacheAble interface {
	Set(key string, value interface{}, timeout time.Duration) error
	Get(key string) interface{}
	Delete(key string) error
	IsExists(key string) bool
	ClearAll() error
	SetOption(option CacheOption)
}
type Instance func() CacheAble

var cachers = map[string]Instance{}

func Register(name string, instance Instance) {
	if instance == nil {
		panic("register a nil instance to cache")
	}
	cachers[name] = instance
}

func NewCache(name string) CacheAble {
	if c, ok := cachers[name]; ok {
		return c()
	} else {
		panic(fmt.Sprintf("cache %s not found", name))
	}
}
