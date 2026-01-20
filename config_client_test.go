package redis

import (
	"context"
	"crypto/tls"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetConfig(t *testing.T) {
	Convey("When a configuration is requested with no options", t, func() {
		var cfg ClientConfig
		ctx := context.Background()
		options, err := cfg.Get(ctx)

		Convey("The redis options are set as the defaults", func() {
			So(err, ShouldBeNil)
			So(options.DB, ShouldEqual, 0)
			So(options.Addr, ShouldEqual, "localhost:6379")
		})
	})

	Convey("When a configuration is requested with options", t, func() {
		expectedDatabase := 10
		expectedAddress := "ons.gov.uk"
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
		}

		cfg := ClientConfig{
			Database:  &expectedDatabase,
			Address:   expectedAddress,
			TLSConfig: tlsConfig,
		}
		ctx := context.Background()
		options, err := cfg.Get(ctx)

		Convey("The redis options are set as the defaults", func() {
			So(err, ShouldBeNil)
			So(options.DB, ShouldEqual, expectedDatabase)
			So(options.Addr, ShouldEqual, expectedAddress)
			So(options.TLSConfig, ShouldEqual, tlsConfig)
		})
	})

	Convey("When a configuration is requested with AWS options", t, func() {
		expectedRegion := "eu-west-2"
		expectedUsername := "test-user"
		expectedService := "memorydb"

		cfg := ClientConfig{
			Region:   expectedRegion,
			Username: expectedUsername,
			Service:  expectedService,
		}
		ctx := context.Background()
		options, err := cfg.Get(ctx)

		Convey("The redis options are set as the defaults with AWS credentials provider", func() {
			So(err, ShouldBeNil)
			So(options.DB, ShouldEqual, 0)
			So(options.Addr, ShouldEqual, "localhost:6379")
			So(options.CredentialsProviderContext, ShouldNotBeNil)
		})
	})

	Convey("When a configuration is requested with invalid AWS options", t, func() {
		cfg := ClientConfig{
			Region: "eu-west-2",
		}
		ctx := context.Background()
		_, err := cfg.Get(ctx)

		Convey("Then an error is returned indicating the invalid configuration", func() {
			So(err, ShouldNotBeNil)
		})
	})
}
