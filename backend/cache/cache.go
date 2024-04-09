package cache

import (
	"context"
	"moria/config"
	"net/url"
	"time"

	"github.com/allegro/bigcache/v3"
)

var cache *bigcache.BigCache
var enabled bool = false

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
	enabled = true
}

func Key(imageID string, queryMap url.Values) string {
	if !enabled {
		return ""
	}

	return hash(getKey(imageID, queryMap))
}

func Get(key string) []byte {
	if !enabled {
		return nil
	}

	if key == "" {
		return nil
	}

	data, err := cache.Get(key)
	if err != nil {
		return nil
	}

	return data
}

func Set(key string, data []byte) {
	if !config.IsCacheEnabled() {
		return
	}

	if key == "" {
		return
	}

	err := cache.Set(key, data)
	if err != nil {
		panic(err)
	}
}

func getKey(imageID string, queryMap url.Values) string {
	if len(queryMap) == 0 {
		return imageID
	}

	return imageID + queryMap.Encode()
}
