package redis

import (
	"context"
	"fmt"
	"testing"

	"github.com/ONSdigital/dis-redis/mocks"

	redis "github.com/redis/go-redis/v9"
	. "github.com/smartystreets/goconvey/convey"
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
			if key == "testKey" {
				cmd.SetVal("testValue")
			} else {
				cmd.SetErr(redis.Nil)
			}
			return cmd
		}

		client := &Client{
			redisClient: mockRedisClient,
		}

		Convey("When the key exists", func() {
			val, err := client.GetValue(context.Background(), "testKey")

			Convey("It should return the correct value and no error", func() {
				So(err, ShouldBeNil)
				So(val, ShouldEqual, "testValue")
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
