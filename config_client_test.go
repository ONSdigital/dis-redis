package redis

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetConfig(t *testing.T) {
	Convey("When a configuration is requested with no options", t, func() {
		var cfg ClientConfig

		options, err := cfg.Get()

		Convey("The redis options are set as the defaults", func() {
			So(err, ShouldBeNil)
			So(options.DB, ShouldEqual, 0)
			So(options.Addr, ShouldEqual, "localhost:6379")
		})
	})

	Convey("When a configuration is requested with options", t, func() {
		expectedDatabase := 10
		expectedAddress := "ons.gov.uk"

		cfg := ClientConfig{
			Database: &expectedDatabase,
			Address:  expectedAddress,
		}

		options, err := cfg.Get()

		Convey("The redis options are set as the defaults", func() {
			So(err, ShouldBeNil)
			So(options.DB, ShouldEqual, expectedDatabase)
			So(options.Addr, ShouldEqual, expectedAddress)
		})
	})
}
