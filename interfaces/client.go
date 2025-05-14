package interfaces

import redis "github.com/redis/go-redis/v9"

//go:generate moq -out ../mocks/go-redis_client.go -pkg mocks . GoRedisClient

// RedisClient is an alias for redis.UniversalClient
type GoRedisClient = redis.UniversalClient
