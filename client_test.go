package redis

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/ONSdigital/dis-redis/mocks"
	"github.com/redis/go-redis/v9"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	TestKey   = "testKey"
	TestValue = "testValue"
)

func TestNewClient(t *testing.T) {
	Convey("When a client is created with no options passed", t, func() {
		ctx := context.Background()

		client, err := NewClient(ctx, &ClientConfig{})

		Convey("The default configuration should be present", func() {
			So(err, ShouldBeNil)

			expectedDatabase := 0
			expectedAddress := "localhost:6379"

			if c, ok := client.redisClient.(*redis.Client); ok {
				actualOptions := c.Options()

				So(actualOptions.DB, ShouldEqual, expectedDatabase)
				So(actualOptions.Addr, ShouldEqual, expectedAddress)
			}
		})
	})

	Convey("When a client is created with a fake redis client", t, func() {
		ctx := context.Background()

		expectedResponse := "pong"

		goRedisClientMock := &mocks.GoRedisClientMock{
			GetFunc: func(ctx context.Context, key string) *redis.StringCmd {
				cmd := redis.NewStringCmd(ctx, "get", key)
				cmd.SetVal(expectedResponse)
				return cmd
			},
		}

		client := NewClientWithCustomClient(ctx, &ClientConfig{}, goRedisClientMock)

		Convey("Then a mock function can be called with a mocked response", func() {
			actualResponse, err := client.redisClient.Get(ctx, "ping").Result()
			So(err, ShouldBeNil)
			So(actualResponse, ShouldEqual, expectedResponse)
		})
	})
}

func TestClient_GetValue(t *testing.T) {
	mockRedisClient := &mocks.GoRedisClientMock{}

	Convey("Given a mock Redis client", t, func() {
		mockRedisClient.GetFunc = func(ctx context.Context, key string) *redis.StringCmd {
			cmd := redis.NewStringCmd(ctx, key)
			if key == TestKey {
				cmd.SetVal(TestValue)
			} else {
				cmd.SetErr(redis.Nil)
			}
			return cmd
		}

		client := &Client{
			redisClient: mockRedisClient,
		}

		Convey("When the key exists", func() {
			val, err := client.GetValue(context.Background(), TestKey)

			Convey("It should return the correct value and no error", func() {
				So(err, ShouldBeNil)
				So(val, ShouldEqual, TestValue)
			})
		})

		Convey("When the key does not exist", func() {
			val, err := client.GetValue(context.Background(), "nonExistingKey")

			Convey("It should return an error with a 'not found' message", func() {
				So(err, ShouldNotBeNil)
				So(val, ShouldBeEmpty)
				So(err.Error(), ShouldContainSubstring, "not found")
			})
		})

		Convey("When Redis returns an error", func() {
			mockRedisClient.GetFunc = func(ctx context.Context, key string) *redis.StringCmd {
				cmd := redis.NewStringCmd(ctx, key)
				cmd.SetErr(fmt.Errorf("connection error"))
				return cmd
			}

			val, err := client.GetValue(context.Background(), "someKey")

			Convey("It should return the correct error message", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "connection error")
				So(val, ShouldBeEmpty)
			})
		})
	})
}

func Test_GetKeyValuePairs(t *testing.T) {
	ctx := context.Background()

	mockRedisClient := &mocks.GoRedisClientMock{}

	// Simulating the cmdable function type
	mockCmdable := func(ctx context.Context, cmd redis.Cmder) error {
		switch c := cmd.(type) {
		case *redis.ScanCmd:
			c.SetVal([]string{"key1", "key2"}, 0)
		case *redis.SliceCmd:
			// Simulate the behavior of MGetCmd
			c.SetVal([]interface{}{"value1", "value2"})
		}
		return nil
	}

	client := &Client{
		redisClient: mockRedisClient,
	}

	// Test case: Successfully retrieve key-value pairs
	Convey("Given a mocked Redis client", t, func() {
		Convey("When retrieving key-value pairs", func() {
			// Set the mock behavior for Scan and MGet
			mockRedisClient.ScanFunc = func(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {
				// Mock ScanCmd using available functions without accessing baseCmd
				cmd := redis.NewScanCmd(ctx, mockCmdable, cursor, match, count) // No need to access baseCmd directly
				cmd.SetVal([]string{"key1", "key2"}, cursor)
				return cmd
			}

			// Mock MGet with the correct conversion of []string to []interface{}
			mockRedisClient.MGetFunc = func(ctx context.Context, keys ...string) *redis.SliceCmd {
				interfaceKeys := make([]interface{}, len(keys))
				for i, key := range keys {
					interfaceKeys[i] = key
				}

				cmd := redis.NewSliceCmd(ctx, interfaceKeys...)
				cmd.SetVal([]interface{}{"value1", "value2"})
				return cmd
			}

			keyValuePairs, err := client.GetKeyValuePairs(ctx, "prefix:*", 10)

			Convey("Then it should return the correct key-value pairs", func() {
				So(err, ShouldBeNil)
				So(keyValuePairs, ShouldResemble, map[string]string{
					"key1": "value1",
					"key2": "value2",
				})
			})

			Convey("Then it should return an empty map when no keys match", func() {
				mockRedisClient.ScanFunc = func(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {
					cmd := redis.NewScanCmd(ctx, mockCmdable, match, count)
					cmd.SetVal([]string{}, 0) // Return no keys
					return cmd
				}
				keyValuePairs, err := client.GetKeyValuePairs(ctx, "prefix:*", 10)
				So(err, ShouldBeNil)
				So(keyValuePairs, ShouldResemble, map[string]string{}) // No keys found
			})

			Convey("Then it should handle cursor pagination correctly", func() {
				// Mock multiple cursors for pagination
				mockRedisClient.ScanFunc = func(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {
					if cursor == 0 {
						cmd := redis.NewScanCmd(ctx, mockCmdable, match, count)
						cmd.SetVal([]string{"key1", "key2"}, 123) // Simulate first page with cursor 123
						return cmd
					}
					cmd := redis.NewScanCmd(ctx, mockCmdable, match, count)
					cmd.SetVal([]string{"key3", "key4"}, 0) // Simulate second page with cursor 0 (end)
					return cmd
				}

				keyValuePairs, err := client.GetKeyValuePairs(ctx, "prefix:*", 10)
				So(err, ShouldBeNil)
				So(keyValuePairs, ShouldResemble, map[string]string{
					"key1": "value1",
					"key2": "value2",
					"key3": "value1",
					"key4": "value2",
				})
			})
		})
	})
}

func TestClient_GetTotalKeys(t *testing.T) {
	mockRedisClient := &mocks.GoRedisClientMock{}
	client := &Client{
		redisClient: mockRedisClient,
	}

	Convey("Given a mocked Redis client", t, func() {
		mockRedisClient.DBSizeFunc = func(ctx context.Context) *redis.IntCmd {
			cmd := redis.NewIntCmd(ctx)
			cmd.SetVal(100) // Mock that there are 100 keys in the Redis DB
			return cmd
		}

		Convey("When calling GetTotalKeys", func() {
			totalKeys, err := client.GetTotalKeys(context.Background())

			Convey("Then it should return the correct total number of keys", func() {
				So(err, ShouldBeNil)
				So(totalKeys, ShouldEqual, 100)
			})
		})
	})

	Convey("Given a mocked Redis client that returns an error", t, func() {
		mockRedisClient.DBSizeFunc = func(ctx context.Context) *redis.IntCmd {
			cmd := redis.NewIntCmd(ctx)
			cmd.SetErr(errors.New("Redis error"))
			return cmd
		}

		Convey("When calling GetTotalKeys", func() {
			totalKeys, err := client.GetTotalKeys(context.Background())

			Convey("Then it should return an error", func() {
				So(err, ShouldNotBeNil)
				So(totalKeys, ShouldEqual, 0)
				So(err.Error(), ShouldContainSubstring, "Redis error")
			})
		})
	})
}

func TestClient_SetValue(t *testing.T) {
	ctx := context.Background()
	mockRedisClient := &mocks.GoRedisClientMock{}
	called := false // Flag to track if Set function was called

	client := &Client{
		redisClient: mockRedisClient,
	}

	Convey("Given a mocked Redis client with expiration set", t, func() {
		Convey("When setting a key-value pair with expiration", func() {
			mockRedisClient.SetFunc = func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
				called = true
				if key == TestKey && value == TestValue && expiration == 10*time.Second {
					cmd := redis.NewStatusCmd(ctx, "set", key)
					cmd.SetVal("OK")
					return cmd
				}
				return redis.NewStatusCmd(ctx, "set", key)
			}

			err := client.SetValue(ctx, TestKey, TestValue, 10*time.Second)

			Convey("Then it should not return an error", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then all expected Redis commands should have been called", func() {
				So(called, ShouldBeTrue)
			})
		})
	})

	Convey("Given a mocked Redis client without expiration", t, func() {
		Convey("When setting a key-value pair without expiration", func() {
			mockRedisClient.SetFunc = func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
				called = true
				if key == TestKey && value == TestValue && expiration == 0 {
					cmd := redis.NewStatusCmd(ctx, "set", key)
					cmd.SetVal("OK")
					return cmd
				}
				return redis.NewStatusCmd(ctx, "set", key)
			}

			err := client.SetValue(ctx, TestKey, TestValue, 0)

			Convey("Then it should not return an error", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then all expected Redis commands should have been called", func() {
				So(called, ShouldBeTrue)
			})
		})
	})
}

func TestClient_DeleteValue(t *testing.T) {
	ctx := context.Background()
	mockRedisClient := &mocks.GoRedisClientMock{}
	called := false // Flag to track if Del function was called
	client := &Client{
		redisClient: mockRedisClient,
	}

	Convey("Given a mocked Redis client", t, func() {
		Convey("When deleting a key that exists", func() {
			mockRedisClient.DelFunc = func(ctx context.Context, key ...string) *redis.IntCmd {
				called = true
				cmd := redis.NewIntCmd(ctx, "del", key)
				cmd.SetVal(1) // Simulate successful deletion (key found)
				return cmd
			}

			err := client.DeleteValue(ctx, "testKey")

			Convey("Then it should not return an error", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the Del function should have been called", func() {
				So(called, ShouldBeTrue)
			})
		})

		Convey("When deleting a key that does not exist", func() {
			mockRedisClient.DelFunc = func(ctx context.Context, key ...string) *redis.IntCmd {
				called = true
				cmd := redis.NewIntCmd(ctx, "del", key)
				cmd.SetVal(0) // Simulate key not found
				return cmd
			}

			// Call DeleteValue function
			err := client.DeleteValue(ctx, "nonExistingKey")

			Convey("Then it should return a 'key not found' error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "key not found")
			})
		})

		Convey("When Redis returns an error", func() {
			mockRedisClient.DelFunc = func(ctx context.Context, key ...string) *redis.IntCmd {
				called = true
				cmd := redis.NewIntCmd(ctx, "del", key)
				cmd.SetErr(fmt.Errorf("connection error"))
				return cmd
			}

			err := client.DeleteValue(ctx, "someKey")

			Convey("Then it should return the correct error message", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "connection error")
			})
		})
	})
}
