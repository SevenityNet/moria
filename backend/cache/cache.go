package cache

import (
	"context"
	"moria/config"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/google/uuid"
)

var cache *bigcache.BigCache

// Initializes the cache.
func Initialize() {
	if !config.IsCacheEnabled() {
		return
	}

	c, err := bigcache.New(context.Background(), bigcache.DefaultConfig(time.Second*time.Duration(config.GetCacheExpiration())))
	if err != nil {
		panic(err)
	}

	cache = c
}

func Get(imageID uuid.UUID) []byte {
	if !config.IsCacheEnabled() {
		return nil
	}

	key := hash(imageID.String())
	if key == "" {
		return nil
	}

	data, err := cache.Get(key)
	if err != nil {
		return nil
	}

	return data
}

func Set(imageID uuid.UUID, data []byte) {
	if !config.IsCacheEnabled() {
		return
	}

	key := hash(imageID.String())
	if key == "" {
		return
	}

	err := cache.Set(key, data)
	if err != nil {
		panic(err)
	}
}
