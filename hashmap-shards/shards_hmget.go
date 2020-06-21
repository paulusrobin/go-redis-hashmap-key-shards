package hashmap_shards

func (impl implementation) HMGet(key string, fields ...string) ([]interface{}, error) {
	var response = make([]interface{}, 0)
	var keys = impl.constructMapKey(key, fields...)
	var dataChannel = make(chan []interface{}, len(keys))

	for k, v := range keys {
		go impl.hmGet(k, dataChannel, v...)
	}

	for i := 0; i < len(keys); i++ {
		select {
		case data := <-dataChannel:
			response = append(response, data...)
			break
		}
	}
	close(dataChannel)
	return response, nil
}

func (impl implementation) hmGet(key string, data chan []interface{}, fields ...string) {
	response, err := impl.HMGet(key, fields...)
	if err != nil {
		response = make([]interface{}, 0)
		impl.log.Errorf("error on hmGet with key: %s and fields %+v", key, fields)
	}
	data <- response
}
