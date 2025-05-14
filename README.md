# dis-redis

Dissemination client library for Redis or similar protocols. This library can be considered a wrapper around go-redis.

## redis package

Includes implementation of a health checker, that reuses the redis client to check requests can be made against the redis server.

### Setup dis-redis client

```golang
import disRedis "github.com/ONSdigital/dis-redis"

...
    redisClient, redisClientErr := disRedis.NewRedisClient(ctx, &disRedis.ClientConfig{
        Address: cfg.redisURL
    }
    if redisClientErr != nil {
        log.Fatal(ctx, "Failed to create dis-redis client", redisClientErr)
    }
...

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

* No further dependencies other than those defined in `go.mod`

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

## License

Copyright © 2025, Office for National Statistics <https://www.ons.gov.uk>

Released under MIT license, see [LICENSE](LICENSE.md) for details.
