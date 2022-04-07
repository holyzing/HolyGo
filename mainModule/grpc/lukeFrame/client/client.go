package client

import (
	"context"
	"fmt"
	"mainModule/grpc/lukeFrame/client/transport"
	"os"
	"strings"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/log"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
)

// Service information constants
const (
	ServiceName  = "luke-service" // Name of service.
	ServiceProxy = "envoy"        // TODO 服务发现, 注册代理 ???
	EndpointsID  = ServiceName    // ID to lookup a service endpoint with.
	maxMsgSize   = 20 * 1024 * 1024
)

// ------------------------------------

type Logger interface {
	// Error logs a message at error priority
	Error(msg string)

	// Infof logs a message at info priority
	Infof(msg string, args ...interface{})
}

type Option struct {
	AccessToken string
	Region      string
	Env         string
	Scheme      string
	Addr        string
	Log         Logger
}

func (o *Option) GetRequestMetadata(ctx context.Context, _ ...string) (map[string]string, error) {
	md := make(map[string]string)
	md["access_token"] = o.AccessToken
	return md, nil
}

func (o *Option) RequireTransportSecurity() bool {
	return false
}

func (o *Option) MergeIn(options ...Option) {
	for _, opt := range options {
		if opt.AccessToken != "" {
			o.AccessToken = opt.AccessToken
		}
		if opt.Region != "" {
			o.Region = strings.ToUpper(opt.Region)
		}
		if opt.Env != "" {
			o.Env = strings.ToUpper(opt.Env)
		}
		if opt.Scheme != "" {
			o.Scheme = strings.ToLower(opt.Scheme)
		}
		if opt.Addr != "" {
			o.Addr = opt.Addr
		}
		if opt.Log != nil {
			o.Log = opt.Log
		}
	}
}

func (o *Option) Init() {
	if o.Addr != "" {
		return
	}

	if o.Scheme == "" {
		o.Scheme = "https"
	}
	if o.Region == "" {
		o.Region = "" // quigon.GetNetworkRegion()
	}
	if o.Env == "" {
		o.Env = "" // quigon.GetEnvironment()
	}

	switch o.Region {
	case "ALPHA":
		switch o.Env {
		case "TESTING":
			if o.Scheme == "grpc" {
				o.Addr = ServiceProxy + ".testing:80"
			} else {
				o.Addr = "/luke-service.Testing"
			}
		case "DEVELOPMENT":
			if o.Scheme == "grpc" {
				o.Addr = ServiceProxy + ".development:80"
			} else {
				o.Addr = "/luke-service.Development"
			}
		default:
			if o.Scheme == "grpc" {
				o.Addr = "localhost:5040"
			} else {
				o.Addr = "localhost:5050"
			}
		}
	case "CHI":
		switch o.Env {
		case "PRODUCTION":
			if o.Scheme == "grpc" {
				o.Addr = ServiceProxy + ".production:80"
			} else {
				o.Addr = "/luke-service.Production"
			}
		case "STAGING":
			if o.Scheme == "grpc" {
				o.Addr = ServiceProxy + ".staging:80"
			} else {
				o.Addr = "/luke-service.Staging"
			}
		default:
			if o.Scheme == "grpc" {
				o.Addr = "localhost:5040"
			} else {
				o.Addr = "localhost:5050"
			}
		}
	case "RHO":
		switch o.Env {
		case "PRODUCTION":
			if o.Scheme == "grpc" {
				o.Addr = ServiceProxy + ".production:80"
			} else {
				o.Addr = "/luke-service.Production"
			}
		case "STAGING":
			if o.Scheme == "grpc" {
				o.Addr = ServiceProxy + ".staging:80"
			} else {
				o.Addr = "/luke-service.Staging"
			}
		default:
			if o.Scheme == "grpc" {
				o.Addr = "localhost:5040"
			} else {
				o.Addr = "localhost:5050"
			}
		}
	}

}

// ------------------------------------

type Luke struct {
	transport.Set
	option Option
	tracer stdopentracing.Tracer
}

// New init client
func New(opts ...Option) *Luke {
	// This is a demonstration client, which supports multiple transports.
	// Your clients will probably just define and stick with 1 transport.
	var (
		svc transport.Set
		err error
	)
	opt := new(Option)
	opt.MergeIn(opts...)
	opt.Init()
	otTracer := NewTracer(opt.Log)

	switch opt.Scheme {
	// case "http", "https":
	// 	svc, err = transport.NewHTTPClient(
	// 		opt.Scheme+"://"+opt.Addr,
	// 		opentracing.ContextToHTTP(otTracer, log.NewNopLogger()),
	// 		transport.SessionHttpHeaders(session, opt.AccessToken),
	// 	)
	case "grpc":
		opts := []grpc.DialOption{
			//lint:file-ignore SA1019 建议的类型不准确
			grpc.WithInsecure(),
			//lint:ignore SA1019 建议的类型不准确
			grpc.WithTimeout(10 * time.Second),

			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(maxMsgSize),
				grpc.MaxCallSendMsgSize(maxMsgSize),
			),
		}
		if opt.AccessToken != "" {
			opts = append(opts, grpc.WithPerRPCCredentials(opt))
		} else {
			println("quigon session")
			// opts = append(opts, grpc.WithPerRPCCredentials(session))
		}
		conn, err := grpc.Dial(opt.Addr, opts...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v", err)
			os.Exit(1)
		}

		svc = transport.NewGRPCClient(conn,
			opentracing.ContextToGRPC(otTracer, log.NewNopLogger()),
			// transport.SessionGRPCMetadata(session),
		)
	default:
		panic("not support")
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	luke := &Luke{Set: svc, tracer: otTracer, option: *opt}
	luke.Set = luke.WrapEndpoints(svc)
	return luke
}

func NewTracer(log Logger) stdopentracing.Tracer {
	var err error
	//var closer io.Closer
	jaegerHost := os.Getenv("JAEGER_AGENT_HOST")
	jaegerPort := os.Getenv("JAEGER_AGENT_PORT")
	if jaegerHost == "" {
		jaegerHost = "10.240.2.81"
	}
	if jaegerPort == "" {
		jaegerPort = "5775"
	}
	// fmt.Println("Jarger Addr:", jaegerHost)
	addr := fmt.Sprintf("%s:%s", jaegerHost, jaegerPort)
	cfg := config.Configuration{
		ServiceName: ServiceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  addr,
		},
	}
	otTracer, _, err := cfg.NewTracer(
	//config.Metrics(metrics.NullFactory),
	)
	if log != nil {
		otTracer, _, err = cfg.NewTracer(
			config.Logger(log),
			//config.Metrics(metrics.NullFactory),
		)
	}
	if err != nil {
		panic(err)
	}
	stdopentracing.SetGlobalTracer(otTracer)
	otTracer.StartSpan("CCCC")
	//defer closer.Close()
	return otTracer
}

// WrapEndpoints accepts the service's entire collection of endpoints, so that a
// set of middlewares can be wrapped around every middleware (e.g., access
// logging and instrumentation), and others wrapped selectively around some
// endpoints and not others (e.g., endpoints requiring authenticated access).
// Note that the final middleware wrapped will be the outermost middleware
// (i.e. applied first)

var labeleMiddlewareWithTracer = func(tracer stdopentracing.Tracer) transport.LabeledMiddleware {
	return func(name string, in endpoint.Endpoint) endpoint.Endpoint {
		return opentracing.TraceClient(tracer, name)(in)
	}
}

func (l *Luke) WrapEndpoints(in transport.Set) transport.Set {
	// tracer middleware
	in.WrapAllLabeledExcept(labeleMiddlewareWithTracer(l.tracer))

	return in
}
