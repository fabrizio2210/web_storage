package rediswrapper

import (
  "context"
  "log"

  "github.com/go-redis/redis/v8"
)


var RedisClient *redis.Client
func ConnectRedis(address string) *redis.Client {
  log.Printf("Connecting to \"%s\" for Redis", address)
  return redis.NewClient(&redis.Options{
    Addr: address,
  })
}

var ctx = context.Background()

func Store(key string, data []byte) error{
  if err := RedisClient.Set(ctx, key, data, 0).Err(); err != nil {
    panic(err)
  }
  return nil
}

func Get(key string) []byte{
  val, err := RedisClient.Get(ctx, key).Result()
  if err != nil {
    return nil
  }
  return []byte(val)
}

func Delete(key string) int64{
  res, _ := RedisClient.Del(ctx, key).Result()
  return res
}
