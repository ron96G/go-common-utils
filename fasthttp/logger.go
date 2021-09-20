package fasthttp

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasttemplate"
)

/*
	Configuraable access logger middleware for fasthttp
	Implementation based on https://github.com/labstack/echo/blob/master/middleware/logger.go
*/

const (
	requestIDHeader = "X-Request-ID"
)

type LoggerConfig struct {
	Format string `yaml:"format"`

	TimeFormat string `yaml:"time_format"`

	Output io.Writer

	template *fasttemplate.Template
	pool     *sync.Pool
	envCache map[string][]byte
}

func LoggerWithConfig(h fasthttp.RequestHandler, config LoggerConfig) fasthttp.RequestHandler {

	config.template = fasttemplate.New(config.Format, "${", "}")
	config.pool = &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 256))
		},
	}
	config.envCache = map[string][]byte{}
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		var err error
		start := time.Now()
		req := ctx.Request
		res := ctx.Response
		h(ctx)
		stop := time.Now()
		buf := config.pool.Get().(*bytes.Buffer)
		buf.Reset()
		defer config.pool.Put(buf)

		if _, err = config.template.ExecuteFunc(buf, func(w io.Writer, tag string) (int, error) {
			switch tag {
			case "time":
				return buf.WriteString(time.Now().Format(config.TimeFormat))
			case "id":
				id := req.Header.Peek(requestIDHeader)
				if id == nil {
					id = res.Header.Peek(requestIDHeader)
				}
				return buf.Write(id)
			case "remote_ip":
				return buf.WriteString(ctx.RemoteAddr().String())
			case "host":
				return buf.Write(req.Host())
			case "uri":
				return buf.Write(req.URI().FullURI())
			case "method":
				return buf.Write(req.Header.Method())
			case "path":
				p := req.URI().Path()
				if p == nil {
					p = []byte("/")
				}
				return buf.Write(p)
			case "protocol":
				return buf.Write(req.Header.Protocol())
			case "referer":
				return buf.Write(req.Header.Referer())
			case "user_agent":
				return buf.Write(req.Header.UserAgent())
			case "status":
				return buf.WriteString(strconv.Itoa(res.StatusCode()))
			case "error":
				if err != nil {
					// Error may contain invalid JSON e.g. `"`
					b, _ := json.Marshal(err.Error())
					b = b[1 : len(b)-1]
					return buf.Write(b)
				}
			case "latency_sec":
				l := stop.Sub(start)
				return buf.WriteString(strconv.FormatFloat(float64(l)/float64(time.Second), 'f', 4, 64))
			case "latency":
				l := stop.Sub(start)
				return buf.WriteString(strconv.FormatInt(int64(l), 10))
			case "latency_human":
				return buf.WriteString(stop.Sub(start).String())
			case "bytes_in":
				cl := req.Header.Peek("Content-Length")
				if cl == nil {
					cl = []byte("0")
				}
				return buf.Write(cl)
			case "bytes_out":
				return buf.WriteString(strconv.FormatInt(int64(res.Header.ContentLength()), 10))
			default:
				switch {
				case strings.HasPrefix(tag, "header:"):
					return buf.Write(req.Header.Peek(tag[7:]))
				case strings.HasPrefix(tag, "env:"):
					val, found := config.envCache[tag[4:]]
					if found {
						return buf.Write(val)
					}
					val = []byte(os.Getenv(tag[4:]))
					config.envCache[tag[4:]] = val
					return buf.Write(val)
				}
			}
			return 0, nil
		}); err != nil {
			return
		}

		_, err = config.Output.Write(buf.Bytes())
	})
}
