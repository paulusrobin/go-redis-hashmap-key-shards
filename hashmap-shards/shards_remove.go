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

	for _, key := range keys {
		err := impl.redis.Remove(key)
		if err != nil {
			errs = append(errs, err.Error())
			impl.log.Error(err)
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}

	return nil
}