/*
 * Copyright (c) 2019.
 */

package cache

import (
	"encoding/gob"
	"io"
	"os"
	"sync"
	"time"
)

var (
	filecache_dir = "cache/.cinyoung/"
)

type fileCacheItem struct {
	ExpireTime time.Time
	Forever    bool
	Value      interface{}
}

type FileCache struct {
	sync.RWMutex
	cacheDir string
}

func (fc *FileCache) Set(key string, value interface{}, timeout time.Duration) error {
	var fdir = fc.cacheDir + "/" + filecache_dir
	_ = os.MkdirAll(fdir, 0666)
	var fname = fdir + key
	item := fileCacheItem{
		Value:      value,
		ExpireTime: time.Now().Add(timeout),
	}
	if timeout == time.Duration(0) {
		item.Forever = true
	}
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	return fc.WriteTo(f, item)
}

func (fc *FileCache) Get(key string) interface{} {
	var fname = fc.cacheDir + "/" + filecache_dir + key
	f, err := os.OpenFile(fname, os.O_RDONLY, 0666)
	if err != nil {
		return nil
	}
	defer f.Close()
	var item = fileCacheItem{}
	err = fc.ReadFrom(f, &item)
	if err != nil {
		return nil
	}
	if item.ExpireTime.Unix() <= time.Now().Unix() && !item.Forever {
		//过期，删除
		defer os.Remove(fname)
		return nil
	}
	return item.Value
}

func (fc *FileCache) Delete(key string) error {
	var fname = fc.cacheDir + "/" + filecache_dir + key
	return os.Remove(fname)
}

func (fc *FileCache) IsExists(key string) bool {
	var fname = fc.cacheDir + "/" + filecache_dir + key
	f, err := os.OpenFile(fname, os.O_RDONLY, 0666)
	if err != nil {
		return false
	}
	defer f.Close()
	var item = fileCacheItem{}
	err = fc.ReadFrom(f, &item)
	if err != nil {
		return false
	}
	if item.ExpireTime.Unix() <= time.Now().Unix() && !item.Forever {
		//过期，删除
		defer os.Remove(fname)
		return false
	}
	return true
}

func (fc *FileCache) ClearAll() error {
	var Dir = fc.cacheDir + "/" + filecache_dir
	return os.RemoveAll(Dir)
}

func (fc *FileCache) SetOption(option CacheOption) {
}

func (FileCache) WriteTo(w io.Writer, item fileCacheItem) error {
	return gob.NewEncoder(w).Encode(item)
}

func (FileCache) ReadFrom(r io.Reader, item *fileCacheItem) error {
	return gob.NewDecoder(r).Decode(item)
}

func NewFileCache() CacheAble {
	return &FileCache{
		cacheDir: os.TempDir(),
	}
}

func init() {
	gob.Register(FileCache{})
	Register("file", NewFileCache)
}
