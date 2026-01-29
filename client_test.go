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
	TestKey     = "testKey"
	TestValue   = "testValue"
	testAddress = "localhost:6379"
)

func TestNewClient(t *testing.T) {
	Convey("When a client is created with no options passed", t, func() {
		ctx := context.Background()

		client, err := NewClient(ctx, &ClientConfig{})

		Convey("The default configuration should be present", func() {
			So(err, ShouldBeNil)

			expectedDatabase := 0
			expectedAddress := testAddress

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

func TestNewClusterClient(t *testing.T) {
	Convey("When a ClusterClient is created with no options passed", t, func() {
		ctx := context.Background()

		client, err := NewClusterClient(ctx, &ClientConfig{})

		Convey("The default configuration should be present", func() {
			So(err, ShouldBeNil)

			expectedDatabase := 0
			expectedAddress := testAddress

			if c, ok := client.redisClient.(*redis.Client); ok {
				actualOptions := c.Options()

				So(actualOptions.DB, ShouldEqual, expectedDatabase)
				So(actualOptions.Addr, ShouldEqual, expectedAddress)
			}
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
				So(err.Error(), ShouldEqual, ErrKeyNotFound.Error())
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

func TestClient_GetKeyValuePairs(t *testing.T) {
	ctx := context.Background()
	match := "prefix:*"
	count := int64(5)

	// Simulating the cmdable function type
	mockCmdable := func(ctx context.Context, cmd redis.Cmder) error {
		return nil
	}

	mockRedisClient := &mocks.GoRedisClientMock{}

	Convey("Given a mock Redis client with paginated data", t, func() {
		// Mock ScanFunc to paginate through 3 pages
		mockRedisClient.ScanFunc = func(ctx context.Context, cursor uint64, pattern string, count int64) *redis.ScanCmd {
			cmd := redis.NewScanCmd(ctx, mockCmdable, cursor, pattern, count)
			switch cursor {
			case 0:
				cmd.SetVal([]string{"key1", "key2", "key3", "key4", "key5"}, 1)
			case 1:
				cmd.SetVal([]string{"key6", "key7", "key8", "key9", "key10"}, 2)
			case 2:
				cmd.SetVal([]string{"key11", "key12", "key13", "key14", "key15"}, 0)
			default:
				cmd.SetVal([]string{}, 0)
			}
			return cmd
		}

		// Mock MGetFunc to return values for given keys
		mockRedisClient.MGetFunc = func(ctx context.Context, keys ...string) *redis.SliceCmd {
			cmd := redis.NewSliceCmd(ctx)
			vals := make([]interface{}, len(keys))
			for i, key := range keys {
				vals[i] = fmt.Sprintf("val_for_%s", key)
			}
			cmd.SetVal(vals)
			return cmd
		}

		client := &Client{
			redisClient: mockRedisClient,
		}

		Convey("When calling GetKeyValuePairs for page 1", func() {
			results, nextCursor, err := client.GetKeyValuePairs(ctx, match, count, 0)

			Convey("It should return first 5 key-value pairs", func() {
				So(err, ShouldBeNil)
				So(nextCursor, ShouldEqual, 1)
				So(results, ShouldResemble, map[string]string{
					"key1": "val_for_key1",
					"key2": "val_for_key2",
					"key3": "val_for_key3",
					"key4": "val_for_key4",
					"key5": "val_for_key5",
				})
			})

			Convey("Then calling GetKeyValuePairs for page 2", func() {
				results2, nextCursor2, err2 := client.GetKeyValuePairs(ctx, match, count, nextCursor)

				Convey("It should return next 5 key-value pairs", func() {
					So(err2, ShouldBeNil)
					So(nextCursor2, ShouldEqual, 2)
					So(results2, ShouldResemble, map[string]string{
						"key6":  "val_for_key6",
						"key7":  "val_for_key7",
						"key8":  "val_for_key8",
						"key9":  "val_for_key9",
						"key10": "val_for_key10",
					})
				})

				Convey("Then calling GetKeyValuePairs for page 3", func() {
					results3, nextCursor3, err3 := client.GetKeyValuePairs(ctx, match, count, nextCursor2)

					Convey("It should return final 5 key-value pairs and cursor 0", func() {
						So(err3, ShouldBeNil)
						So(nextCursor3, ShouldEqual, 0)
						So(results3, ShouldResemble, map[string]string{
							"key11": "val_for_key11",
							"key12": "val_for_key12",
							"key13": "val_for_key13",
							"key14": "val_for_key14",
							"key15": "val_for_key15",
						})
					})
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
