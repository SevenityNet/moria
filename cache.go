package main

import (
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
)

var (
	AUTHCACHE *bigcache.BigCache
	FILECACHE *bigcache.BigCache
)

func initCache() {
	var err error
	AUTHCACHE, err = bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
	if err != nil {
		panic(err)
	}

	FILECACHE, err = bigcache.New(context.Background(), bigcache.DefaultConfig(1*time.Hour))
	if err != nil {
		panic(err)
	}
}

// caches a file, return error if saving to cache fails and forwards bigcache error
func cacheFile(filename string, file []byte) error {
	err := FILECACHE.Set(filename, file)
	if err != nil {
		return err
	}
	return nil
}

// Return cached file if it exists in cache, otherwise returns ErrFileNotFound. If an error happens in retrieval it forwards bigcache error
func getCachedFileIfExists(filename string) ([]byte, error) {
	entry, err := FILECACHE.Get(filename)
	if err == bigcache.ErrEntryNotFound {
		return nil, ErrFileNotFound
	} else if err != nil {
		return nil, err
	}

	return entry, nil
}

// removes a cached file, returns error if file is not found in cache or if an error happens in removal
func removeCachedFile(filename string) error {
	err := FILECACHE.Delete(filename)
	if err != nil {
		return err
	}
	return nil
}
