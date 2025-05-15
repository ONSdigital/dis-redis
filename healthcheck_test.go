package redis

import (
	"context"
	"errors"
	"testing"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"

	"github.com/ONSdigital/dis-redis/mocks"
	redis "github.com/redis/go-redis/v9"
	. "github.com/smartystreets/goconvey/convey"
)

func TestChecker(t *testing.T) {
	Convey("Given that Redis is healthy", t, func() {
		ctx := context.Background()

		mockRedisClient := &mocks.GoRedisClientMock{
			PingFunc: func(ctx context.Context) *redis.StatusCmd {
				cmd := redis.NewStatusCmd(ctx)
				cmd.SetVal("pong")
				return cmd
			},
		}

		client := NewClientWithCustomClient(ctx, &ClientConfig{}, mockRedisClient)
		checkState := health.NewCheckState("dis-redis-test")

		Convey("Checker updates the CheckState to a successful state", func() {
			client.Checker(context.Background(), checkState)

			So(len(mockRedisClient.PingCalls()), ShouldEqual, 1)
			So(checkState.Status(), ShouldEqual, health.StatusOK)
			So(checkState.Message(), ShouldEqual, MsgHealthy)
			So(checkState.StatusCode(), ShouldEqual, 200)
		})
	})

	Convey("Given that Redis is unhealthy", t, func() {
		ctx := context.Background()

		mockRedisClient := &mocks.GoRedisClientMock{
			PingFunc: func(ctx context.Context) *redis.StatusCmd {
				cmd := redis.NewStatusCmd(ctx)
				cmd.SetErr(errors.New("failed connection"))
				return cmd
			},
		}

		client := NewClientWithCustomClient(ctx, &ClientConfig{}, mockRedisClient)
		checkState := health.NewCheckState("dis-redis-test")

		Convey("Checker updates the CheckState to an unsuccessful state", func() {
			client.Checker(context.Background(), checkState)

			So(len(mockRedisClient.PingCalls()), ShouldEqual, 1)
			So(checkState.Status(), ShouldEqual, health.StatusCritical)
			So(checkState.Message(), ShouldEqual, ErrorFailedConnection.Error())
			So(checkState.StatusCode(), ShouldEqual, 500)
		})
	})
}
