# go-redis-hashmap-key-shards
golang redis hashmap using key shards

### Installation
```
~ $ go get -u github.com/paulusrobin/go-redis-hashmap-key-shards
``` 

#### Example
```go
var redis RedisHashMap
var log Log
var shards uint32
cache, err := hashmap_shards.NewHashMapShards(redis, log, shards)
if err != nil {
    log.Error(err)
    panic(err)
}

_ = cache.HMSetWithExpiration("example", map[string]int{
    "hello": 1,
}, time.Minute)

_ = cache.HSetWithExpiration("example", "world", 2, time.Minute)

if val, err := cache.HMGet("example", "hello", "world"); err == nil {
    fmt.Println(val)
}

if val, err := cache.HGetAll("example"); err == nil {
    fmt.Println(val)
}
```

#### Options
```
# interface RedisHashMap, that implements: 
# - HMSetWithExpiration(key string, value map[string]interface{}, ttl time.Duration) error
# - HMSet(key string, value map[string]interface{}) error
# - HSetWithExpiration(key, field string, value interface{}, ttl time.Duration) error
# - HSet(key, field string, value interface{}) error
# - HMGet(key string, fields ...string) ([]interface{}, error)
# - HGetAll(key string) (map[string]string, error)
# - Remove(key string) error
# - Keys(pattern string) ([]string, error)
Redis       RedisHashMap (interface)

# interface Log, that implements: 
# - Info(...interface{})
# - Infof(string, ...interface{})
# - Debug(...interface{})
# - Debugf(string, ...interface{})
# - Error(...interface{})
# - Errorf(string, ...interface{})
# - Warning(...interface{})
# - Warningf(string, ...interface{})
Log         Log

# Number of Shards use for sharding the key 
Shards      uint32
```

### Performance Benchmarks (Coming Soon)
```
goos: darwin
goarch: amd64
```