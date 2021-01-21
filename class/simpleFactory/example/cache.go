package main

import "errors"

type Cache interface {
	Set(key, value string)
	Get(key string) string
}

type RedisCache struct {
	data map[string]string
}

func (redis *RedisCache) Set(key, value string) {
	redis.data[key] = value
}

func (redis *RedisCache) Get(key string) string {
	return redis.data[key]
}

type MemCache struct {
	data map[string]string
}

func (memCache *MemCache) Set(key, value string) {
	memCache.data[key] = value
}

func (memCache *MemCache) Get(key string) string {
	return memCache.data[key]
}

type cacheType int

const (
	redisCache cacheType = iota
	memCache
)

type cacheFactory struct {
}

func (cf *cacheFactory) Create(cacheType cacheType) (Cache, error) {
	switch cacheType {
	case redisCache:
		return &RedisCache{
			data: map[string]string{},
		}, nil
	case memCache:
		return &MemCache{
			data: map[string]string{},
		}, nil
	default:
		return nil, errors.New("error cache type")
	}
}
