package hashmap_shards

type sliceInterfaceData struct {
	data []interface{}
	err  error
}

func (impl implementation) HMGet(key string, fields ...string) ([]interface{}, error) {
	var response = make([]interface{}, 0)
	var err error = nil
	var keys = impl.constructMapKey(key, fields...)
	var dataChannel = make(chan sliceInterfaceData, len(keys))

	for k, v := range keys {
		go impl.hmGet(k, dataChannel, v...)
	}

	for i := 0; i < len(keys); i++ {
		select {
		case x := <-dataChannel:
			if x.err == nil {
				response = append(response, x.data...)
			} else {
				err = x.err
			}
			break
		}
	}
	close(dataChannel)
	return response, err
}

func (impl implementation) hmGet(key string, data chan sliceInterfaceData, fields ...string) {
	response, err := impl.HMGet(key, fields...)
	data <- sliceInterfaceData{
		data: response,
		err:  err,
	}
}
