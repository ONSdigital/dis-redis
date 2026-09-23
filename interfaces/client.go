package interfaces

import "github.com/redis/go-redis/v9"

//go:generate moq -out ../mocks/go-redis_client.go -pkg mocks . GoRedisClient
//go:generate moq -out ../mocks/go-redis_pipeliner.go -pkg mocks . GoRedisPipeliner

// RedisClient is an alias for redis.UniversalClient
type GoRedisClient = redis.UniversalClient

// GoRedisPipeliner is an alias for redis.Pipeliner.
type GoRedisPipeliner = redis.Pipeliner
