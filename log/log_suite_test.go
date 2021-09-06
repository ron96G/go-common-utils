package log_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ron96G/go-common-utils/log"
)

func TestLog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Log Suite")
}

var _ = Describe("Log", func() {
	defer GinkgoRecover()
	buf := bytes.NewBuffer(nil)

	Describe("Get a new logger", func() {

		log.Warn("hello world")

		log.Configure("error", buf)

		logger := log.New("testlogger")

		logger.Warn("hello world")

		fmt.Println(buf.String())
	})

	buf.Reset()

	Describe("Configure new logger", func() {

		log.Configure("debug", buf, "foo", "bar")

		log.Infof("Test %s", "this")

		fmt.Println(buf.String())
	})

	buf.Reset()

	Describe("Configure loglevel", func() {
		log.Reset()
		log.Configure("error", buf)

		logger := log.New("testlogger")
		log.Warn("Warn")
		logger.Warn("Warn")
		logger.Error("Error")

		fmt.Println(buf.String())
	})

	buf.Reset()

	Describe("Unknown loglevel", func() {
		log.Reset()
		log.Configure("unknown", buf)

		logger := log.New("testlogger")
		log.Warn("hello world")
		logger.Warn("hello world")

		fmt.Println(buf.String())
	})

	Describe("Configure log with params", func() {
		log.Reset()
		log.Configure("info", buf)

		logger := log.New("testlogger", "foo", "bar")
		logger.Warn("hello world")

		fmt.Println(buf.String())
	})

	Describe("Context logger", func() {
		log.Reset()
		log.Configure("info", buf)
		logger := log.New("testlogger", "foo", "bar")

		ctx := context.Background()
		ctxLog := log.ToContext(ctx, logger, "ctx", "t")

		LoggerFromCtx := log.FromContext(ctxLog)

		LoggerFromCtx.Warn("hello from context")
		LoggerFromCtx.Warn("hello from context")

		fmt.Println(buf.String())
	})
})
