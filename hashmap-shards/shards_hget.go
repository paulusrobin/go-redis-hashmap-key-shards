package hashmap_shards

import "fmt"

type mapInterfaceData struct {
	data map[string]string
	err  error
}

func (impl implementation) HGetAll(key string) (map[string]string, error) {
	var response = make(map[string]string)
	var err error = nil
	var dataChannel = make(chan mapInterfaceData, impl.shards)

	for i := uint32(0); i < impl.shards; i++ {
		keyShard := fmt.Sprintf("%s-%d", key, i)
		go impl.hGetAll(keyShard, dataChannel)
	}

	for i := uint32(0); i < impl.shards; i++ {
		select {
		case x := <-dataChannel:
			if x.err == nil {
				for key, value := range x.data {
					response[key] = value
				}
			} else {
				err = x.err
			}
			break
		}
	}
	close(dataChannel)
	return response, err
}

func (impl implementation) hGetAll(key string, data chan mapInterfaceData) {
	response, err := impl.redis.HGetAll(key)
	data <- mapInterfaceData{
		data: response,
		err:  err,
	}
}
