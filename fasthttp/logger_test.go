package fasthttp_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	log_fttp "github.com/ron96G/go-common-utils/fasthttp"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello there")
}

func fastHTTPErrorHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(http.StatusInternalServerError)
}

func StartNewServer(handler fasthttp.RequestHandler) (client *http.Client) {
	ln := fasthttputil.NewInmemoryListener()

	go func() {
		if err := fasthttp.Serve(ln, handler); err != nil {
			panic(err)
		}
	}()
	client = &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return ln.Dial()
			},
		},
		Timeout: time.Second,
	}

	return
}

func TestLoggerWithConfig_Env(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	expVal := "TEST123"
	os.Setenv("FOOBAR", expVal)
	handlerWithLogging := log_fttp.LoggerWithConfig(fastHTTPHandler, log_fttp.LoggerConfig{
		Format:     `"foobar":"${env:FOOBAR}", "hello":"${env:WHAT}"`,
		Output:     buf,
		TimeFormat: time.RFC3339,
	})

	client := StartNewServer(handlerWithLogging)
	_, err := client.Get("http://localhost:8080/irgendwas")
	if err != nil {
		t.Error(err)
	}

	content, _ := ioutil.ReadAll(buf)
	if !bytes.Contains(content, []byte(`"foobar":"TEST123"`)) {
		t.Errorf("Env not set. Expected %s. Found %s", expVal, string(content))
	}
	if !bytes.Contains(content, []byte(`"hello":"-"`)) {
		t.Errorf("Env not set. Expected %s. Found %s", expVal, string(content))
	}
}

func TestLoggerWithSkipperNoError(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	handlerWithLogging := log_fttp.LoggerWithConfig(fastHTTPHandler, log_fttp.LoggerConfig{
		Format:     `"ACCESSLOG"`,
		Output:     buf,
		TimeFormat: time.RFC3339,
		Skipper:    func(ctx *fasthttp.RequestCtx) bool { return ctx.Response.StatusCode() < 400 },
	})

	client := StartNewServer(handlerWithLogging)
	_, err := client.Get("http://localhost:8080/irgendwas")
	if err != nil {
		t.Error(err)
	}

	content, _ := ioutil.ReadAll(buf)
	if bytes.Contains(content, []byte(`"ACCESSLOG"`)) {
		t.Errorf("Skipper should have worked but did not!")
	}
}

func TestLoggerWithSkipperError(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	handlerWithLogging := log_fttp.LoggerWithConfig(fastHTTPErrorHandler, log_fttp.LoggerConfig{
		Format:     `"ACCESSLOG"`,
		Output:     buf,
		TimeFormat: time.RFC3339,
		Skipper:    func(ctx *fasthttp.RequestCtx) bool { return ctx.Response.StatusCode() < 400 },
	})

	client := StartNewServer(handlerWithLogging)
	_, err := client.Get("http://localhost:8080/irgendwas")
	if err != nil {
		t.Error(err)
	}

	content, _ := ioutil.ReadAll(buf)
	if !bytes.Contains(content, []byte(`"ACCESSLOG"`)) {
		t.Errorf("Skipper should have worked but did not!")
	}
}
