package hashmap_shards

import (
	"fmt"
	"time"
)

func (impl implementation) HSetWithExpiration(key, field string, value interface{}, ttl time.Duration) error {
	keyShards := fmt.Sprintf("%s-%d", key, impl.hash(field))
	return impl.redis.HSetWithExpiration(keyShards, field, value, ttl)
}

func (impl implementation) HSet(key, field string, value interface{}) error {
	keyShards := fmt.Sprintf("%s-%d", key, impl.hash(field))
	return impl.redis.HSet(keyShards, field, value)
}