package hashmap_shards

import (
	"time"
)

func (impl implementation) HMSetWithExpiration(key string, value map[string]interface{}, ttl time.Duration) error {
	var values = impl.constructMapValue(key, value)
	var dataChannel = make(chan error, len(values))

	for k, v := range values {
		go impl.hmSetWithExpiration(k, v, ttl, dataChannel)
	}

	err := impl.catchError(len(values), dataChannel)
	close(dataChannel)
	return err
}

func (impl implementation) HMSet(key string, value map[string]interface{}) error {
	var values = impl.constructMapValue(key, value)
	var dataChannel = make(chan error, len(values))

	for k, v := range values {
		go impl.hmSet(k, v, dataChannel)
	}

	err := impl.catchError(len(values), dataChannel)
	close(dataChannel)
	return err
}

func (impl implementation) hmSetWithExpiration(key string, value map[string]interface{}, ttl time.Duration, data chan error) {
	data <- impl.redis.HMSetWithExpiration(key, value, ttl)
}

func (impl implementation) hmSet(key string, value map[string]interface{}, data chan error) {
	data <- impl.redis.HMSet(key, value)
}
