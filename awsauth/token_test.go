package awsauth

import (
	"context"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	envAccessKeyID     = "AWS_ACCESS_KEY_ID"
	envSecretAccessKey = "AWS_SECRET_ACCESS_KEY"

	testAccessKey = "TEST_ACCESS_KEY"
	testSecretKey = "TEST_SECRET_KEY"

	testRegion = "eu-west-2"
	testHost   = "example.memorydb.us-east-1.amazonaws.com:6379"
)

func TestNewTokenGenerator(t *testing.T) {
	Convey("Given an access and secret key are set in the environment", t, func() {
		err := setTestAWSCredentialsEnvironment()
		So(err, ShouldBeNil)

		ctx := context.Background()

		Convey("When NewTokenGenerator is called", func() {
			tokenGen, err := NewTokenGenerator(
				ctx,
				testRegion,
				testHost,
			)
			Convey("Then a TokenGenerator is returned without error", func() {
				So(err, ShouldBeNil)
				So(tokenGen, ShouldNotBeNil)
			})
		})
	})
}

func TestGenerate(t *testing.T) {
	Convey("Given a valid TokenGenerator", t, func() {
		err := setTestAWSCredentialsEnvironment()
		So(err, ShouldBeNil)

		ctx := context.Background()

		tokenGen, err := NewTokenGenerator(
			ctx,
			testRegion,
			testHost,
		)
		So(err, ShouldBeNil)
		So(tokenGen, ShouldNotBeNil)

		Convey("When Generate is called", func() {
			token, err := tokenGen.Generate(ctx)

			Convey("Then a token is returned without error", func() {
				So(err, ShouldBeNil)
				So(token, ShouldNotBeEmpty)
				So(token, ShouldContainSubstring, "AWS4-HMAC-SHA256")
				So(token, ShouldContainSubstring, "Credential=")
				So(token, ShouldContainSubstring, "SignedHeaders=")
				So(token, ShouldContainSubstring, "Signature=")
			})
		})
	})

	Convey("Given a TokenGenerator with no credentials", t, func() {
		err := unsetTestAWSCredentialsEnvironment()
		So(err, ShouldBeNil)

		ctx := context.Background()

		tokenGen, err := NewTokenGenerator(
			ctx,
			testRegion,
			testHost,
		)
		So(err, ShouldBeNil)
		So(tokenGen, ShouldNotBeNil)

		Convey("When Generate is called", func() {
			token, err := tokenGen.Generate(ctx)

			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
				So(token, ShouldBeEmpty)
			})
		})
	})
}

func setTestAWSCredentialsEnvironment() error {
	err := os.Setenv(envAccessKeyID, testAccessKey)
	if err != nil {
		return err
	}
	err = os.Setenv(envSecretAccessKey, testSecretKey)
	if err != nil {
		return err
	}
	return nil
}

func unsetTestAWSCredentialsEnvironment() error {
	err := os.Unsetenv(envAccessKeyID)
	if err != nil {
		return err
	}
	err = os.Unsetenv(envSecretAccessKey)
	if err != nil {
		return err
	}
	return nil
}
