package hashmap_shards

func (impl implementation) Keys(pattern string) ([]string, error) {
	return impl.redis.Keys(pattern)
}