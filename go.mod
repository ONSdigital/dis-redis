module github.com/ONSdigital/dis-redis

go 1.24

require (
	github.com/ONSdigital/dp-healthcheck v1.6.4
	github.com/redis/go-redis/v9 v9.8.0
	github.com/smartystreets/goconvey v1.8.1
)

require (
	github.com/ONSdigital/dp-api-clients-go/v2 v2.266.0 // indirect
	github.com/ONSdigital/dp-net/v3 v3.3.0 // indirect
	github.com/ONSdigital/log.go/v2 v2.4.5 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/smarty/assertions v1.16.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.36.0 // indirect
	go.opentelemetry.io/otel/metric v1.36.0 // indirect
	go.opentelemetry.io/otel/trace v1.36.0 // indirect
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
)

/// [CVE-2022-27191] CWE-noinfo: Allows an attacker to crash a server in certain circumstances involving AddHostKey
tool golang.org/x/crypto
