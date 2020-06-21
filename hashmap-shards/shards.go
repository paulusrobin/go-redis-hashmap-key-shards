package hashmap_shards

import (
	"fmt"
	"github.com/pkg/errors"
	"hash/fnv"
	"time"
)

type (
	RedisHashMap interface {
		HMSetWithExpiration(key string, value map[string]interface{}, ttl time.Duration) error
		HMSet(key string, value map[string]interface{}) error
		HSetWithExpiration(key, field string, value interface{}, ttl time.Duration) error
		HSet(key, field string, value interface{}) error
		HMGet(key string, fields ...string) ([]interface{}, error)
		HGetAll(key string) (map[string]string, error)
	}
	implementation struct {
		redis  RedisHashMap
		log    Log
		shards uint32
	}
)

func NewHashMapShards(redis RedisHashMap, log Log, shards uint32) (RedisHashMap, error) {
	if shards < 1 {
		return nil, errors.New("minimum value key shards is 1")
	}

	if redis == nil {
		return nil, errors.New("redis must not nil")
	}

	if log == nil {
		return nil, errors.New("log must not nil")
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

func (impl implementation) constructMapValue(key string, value map[string]interface{}) map[string]map[string]interface{} {
	var values = make(map[string]map[string]interface{})
	for k, v := range value {
		keyShards := fmt.Sprintf("%s-%d", key, impl.hash(k))
		if val, ok := values[keyShards]; !ok {
			val = map[string]interface{}{k: v}
			values[keyShards] = val
		} else {
			val[k] = v
			values[keyShards] = val
		}
	}
	return values
}

func (impl implementation) constructMapKey(key string, fields ...string) map[string][]string {
	var response = make(map[string][]string)
	for _, field := range fields {
		keyShards := fmt.Sprintf("%s-%d", key, impl.hash(field))
		if val, ok := response[keyShards]; !ok {
			response[keyShards] = []string{field}
		} else {
			response[keyShards] = append(val, field)
		}
	}
	return response
}

func (impl implementation) catchError(len int, dataChannel chan error) error {
	var err error
	for i := 0; i < len; i++ {
		select {
		case e := <-dataChannel:
			if e != nil {
				err = e
			}
			break
		}
	}
	return err
}