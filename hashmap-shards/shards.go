package hashmap_shards

import (
	"github.com/pkg/errors"
	"hash/fnv"
	"time"
)

type (
	HashMapShards interface {
		HMSetWithExpiration(key string, value map[string]interface{}, ttl time.Duration) error
		HMSet(key string, value map[string]interface{}) error
		HSetWithExpiration(key, field string, value interface{}, ttl time.Duration) error
		HSet(key, field string, value interface{}) error
		HMGet(key string, fields ...string) ([]interface{}, error)
		HGetAll(key string) (map[string]string, error)
	}
	implementation struct {
		redis  Redis
		log    Log
		shards uint32
	}
)

func NewHashMapShards(redis Redis, log Log, shards uint32) (HashMapShards, error) {
	if shards < 1 {
		return nil, errors.New("minimum value key shards is 1")
	}
	return &implementation{
		redis:  redis,
		log:    log,
		shards: shards,
	}, nil
}

func (impl implementation) hash(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32() % impl.shards
}