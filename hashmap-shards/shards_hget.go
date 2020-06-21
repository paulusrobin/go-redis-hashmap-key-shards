package hashmap_shards

import "fmt"

func (impl implementation) HGetAll(key string) (map[string]string, error) {
	var response = make(map[string]string)
	var dataChannel = make(chan map[string]string, impl.shards)

	for i := uint32(0); i < impl.shards; i++ {
		keyShard := fmt.Sprintf("%s-%d", key, i)
		go impl.hGetAll(keyShard, dataChannel)
	}

	for i := uint32(0); i < impl.shards; i++ {
		select {
		case x := <-dataChannel:
			for key, value := range x {
				response[key] = value
			}
			break
		}
	}

	return response, nil
}

func (impl implementation) hGetAll(key string, data chan map[string]string) {
	response, err := impl.redis.HGetAll(key)
	if err != nil {
		response = nil
		impl.log.Errorf("error on hGetAll with key: %s", key)
	}
	data <- response
}
