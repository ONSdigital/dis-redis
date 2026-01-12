# dis-redis

Dissemination client library for Redis or similar protocols. This library can be considered a wrapper around go-redis.

## redis package

Includes implementation of a health checker, that reuses the redis client to check requests can be made against the redis server.

### Setup dis-redis client

```golang
import disRedis "github.com/ONSdigital/dis-redis"

...
    redisClient, redisClientErr := disRedis.NewClient(ctx, &disRedis.ClientConfig{
        Address: cfg.redisURL
    }
    if redisClientErr != nil {
        log.Fatal(ctx, "Failed to create dis-redis client", redisClientErr)
    }
...

```

### Health checker

Using dis-redis checker function currently performs a PING request against redis.

The healthcheck will only succeed if the request can be performend and the server responds with a PONG.

Instantiate a dis-redis client

```golang
import disRedis "github.com/ONSdigital/dis-redis"

...
    cli := disRedis.NewClient(ctx, clientConfig)
...
```

Call healthchecker with `cli.Checker(context.Background())` and this will return a check object like so:

```json
{
    "name": "string",
    "status": "string",
    "message": "string",
    "status_code": "int",
    "last_checked": "ISO8601 - UTC date time",
    "last_success": "ISO8601 - UTC date time",
    "last_failure": "ISO8601 - UTC date time"
}
```

## Getting started

Tests and static checks are run via:

```sh
    make audit
    make test
    make lint
    make build
```

* Run `make help` to see full list of make targets

### Dependencies

#### Tools

We use `dis-vulncheck` to do Go auditing, which you will [need to install](https://github.com/ONSdigital/dis-vulncheck).

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

## License

Copyright © 2025, Office for National Statistics <https://www.ons.gov.uk>

Released under MIT license, see [LICENSE](LICENSE.md) for details.
