package hashmap_shards

import (
	"errors"
	"strings"
)

func (impl implementation) Remove(key string) error {
	var errs []string
	keys, err := impl.redis.Keys(key+"*")
	if err != nil {
		return err
	}

	var dataChannel = make(chan error, len(keys))

	for _, key := range keys {
		go impl.remove(key, dataChannel)
	}

	for i := 0; i < len(keys); i++ {
		select {
		case x := <-dataChannel:
			if x != nil {
				errs = append(errs, x.Error())
			}
			break
		}
	}
	close(dataChannel)

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}

	return nil
}

func (impl implementation) remove(key string, data chan error) {
	data <- impl.redis.Remove(key)
}