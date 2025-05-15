package redis

import (
	"context"
	"errors"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
)

const MsgHealthy = "redis is healthy"

var (
	ErrorFailedConnection = errors.New("couldn't connect to redis")
)

// Checker executes all healthchecks and then updates the health state
func (cli *Client) Checker(ctx context.Context, state *health.CheckState) error {
	if state == nil {
		state = &health.CheckState{}
	}

	statusCode, err := cli.healthcheck(ctx)
	if err != nil {
		if updateErr := state.Update(health.StatusCritical, err.Error(), statusCode); updateErr != nil {
			return updateErr
		}

		return nil
	}

	if updateErr := state.Update(health.StatusOK, MsgHealthy, statusCode); updateErr != nil {
		return updateErr
	}

	return nil
}

// healthcheck calls redis to check its health status. This call implements only the logic,
// without providing the Check object, and it's aimed for internal use.
func (cli *Client) healthcheck(ctx context.Context) (code int, err error) {
	err = cli.redisClient.Ping(ctx).Err()
	if err != nil {
		return 500, ErrorFailedConnection
	}

	return 200, nil
}
