package client

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/holyzing/HolyGo/go-kit/luke/client/transport"

	"github.com/uber/jaeger-client-go"

	"github.com/go-kit/kit/endpoint"
	gokitOt "github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/log"
	otg "github.com/opentracing/opentracing-go"
	jgc "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
)

// Service information constants
const (
	ServiceName  = "luke-service" // Name of service.
	ServiceProxy = "envoy"        // TODO 服务发现, 注册代理 ???
	// EndpointsID  = ServiceName    // ID to lookup a service endpoint with.
	maxMsgSize = 20 * 1024 * 1024
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

func (o *Option) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
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
	tracer otg.Tracer

	tracerCloser io.Closer
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
	otTracer, tracerCloser := NewTracer(opt.Log)

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
			// https://github.com/grpc-ecosystem/grpc-opentracing
			// grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(stdopentracing.GlobalTracer())),
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
			fmt.Println(fmt.Fprintf(os.Stderr, "error: %v", err))
			os.Exit(1)
		}

		svc = transport.NewGRPCClient(conn,
			gokitOt.ContextToGRPC(otTracer, log.NewNopLogger()),
			// transport.SessionGRPCMetadata(session),
		)
	default:
		panic("not support")
	}
	if err != nil {
		fmt.Println(fmt.Fprintf(os.Stderr, "error: %v\n", err))
		os.Exit(1)
	}

	luke := &Luke{Set: svc, tracer: otTracer, option: *opt, tracerCloser: tracerCloser}
	luke.Set = luke.WrapEndpoints(svc)
	return luke
}

func NewTracer(log Logger) (otg.Tracer, io.Closer) {
	var err error
	//var closer io.Closer
	jaegerHost := os.Getenv("JAEGER_AGENT_HOST")
	jaegerPort := os.Getenv("JAEGER_AGENT_PORT")
	if jaegerHost == "" {
		jaegerHost = "10.240.2.81"
	}
	if jaegerPort == "" {
		jaegerPort = "6831"
	}
	// fmt.Println("Jarger Addr:", jaegerHost)
	addr := fmt.Sprintf("%s:%s", jaegerHost, jaegerPort)
	cfg := jgc.Configuration{
		ServiceName: ServiceName,
		Sampler: &jgc.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jgc.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  addr,
		},
	}
	if log == nil {
		log = jaeger.StdLogger
	}
	otTracer, tracerCloser, err := cfg.NewTracer(
		jgc.Logger(log),
	// jgc.Metrics(metrics.NullFactory),
	)
	if err != nil {
		panic(err)
	}
	otg.SetGlobalTracer(otTracer)
	// defer closer.Close()

	return otTracer, tracerCloser
}

// WrapEndpoints accepts the service's entire collection of endpoints, so that a
// set of middlewares can be wrapped around every middleware (e.g., access
// logging and instrumentation), and others wrapped selectively around some
// endpoints and not others (e.g., endpoints requiring authenticated access).
// Note that the final middleware wrapped will be the outermost middleware
// (i.e. applied first)

var labeleMiddlewareWithTracer = func(tracer otg.Tracer) transport.LabeledMiddleware {
	return func(name string, in endpoint.Endpoint) endpoint.Endpoint {
		return gokitOt.TraceClient(tracer, name)(in)
	}
}

func (l *Luke) WrapEndpoints(in transport.Set) transport.Set {
	// tracer middleware
	in.WrapAllLabeledExcept(labeleMiddlewareWithTracer(l.tracer))

	return in
}
