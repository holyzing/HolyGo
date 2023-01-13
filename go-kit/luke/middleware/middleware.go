package middleware

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	gokitOt "github.com/go-kit/kit/tracing/opentracing"
	otg "github.com/opentracing/opentracing-go"

	pb "github.com/holyzing/HolyGo/go-kit/luke/api/v1"

	"github.com/go-kit/kit/endpoint"
	"github.com/segmentio/ksuid"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type LabeledMiddleware func(string, endpoint.Endpoint) endpoint.Endpoint

var Auth0Middleware endpoint.Middleware = func(ep endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		println("---- middlewares Auth0Middleware ----")
		accessToken, err := parseAccessToken(ctx)
		if err != nil {
			return &pb.LukeResponse{
				Retcode: pb.Retcode_UNAUTHORIZED,
				Error:   &pb.ErrMsg{Message: err.Error()},
			}, nil
		}
		fmt.Println(accessToken)
		return ep(ctx, request)
	}
}

func parseAccessToken(ctx context.Context) (string, error) {
	println("---- middlewares parseAccessToken ----")
	const prefix = "Bearer "
	var accessToken string

	// get access_token from md -- grpc
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if val, ok := md["access_token"]; ok {
			tmp := strings.TrimPrefix(val[0], "Bearer ")
			accessToken = strings.TrimPrefix(tmp, "bearer ")
			return accessToken, nil
		} else {
			return accessToken, fmt.Errorf("no access token in context")
		}
	}
	// get access_token from headers -- http
	bearer, ok := ctx.Value("authorization").(string)
	if !ok {
		return accessToken, fmt.Errorf("no metadata in context")
	}
	// Case insensitive prefix match. See Issue 22736.
	if len(bearer) < len(prefix) || !strings.EqualFold(bearer[:len(prefix)], prefix) {
		return accessToken, fmt.Errorf("no access token in context")
	}
	accessToken = bearer[len(prefix):]
	return accessToken, nil
}

var LoggingMiddleware LabeledMiddleware = func(name string, endpoint endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		println("---- middlewares LoggingMiddleware ----")
		start := time.Now()
		transport := ctx.Value("transport")

		var sourceIP string

		if transport != nil && transport.(string) == "HTTPJSON" {
			if xForwardedFor := ctx.Value("x-forwarded-for"); xForwardedFor != nil {
				frags := strings.Split(xForwardedFor.(string), ",")
				if len(frags) > 0 {
					sourceIP = strings.TrimSpace(frags[0])
				}
			} else if xRealIP := ctx.Value("x-real-ip"); xRealIP != nil {
				sourceIP = strings.TrimSpace(xRealIP.(string))
			} else if ip, _, err := net.SplitHostPort(strings.TrimSpace(ctx.Value("http-remote-addr").(string))); err == nil {
				sourceIP = ip
			}
		} else {
			pr, ok := peer.FromContext(ctx)
			if !ok {
				return nil, fmt.Errorf("invoke FromContext() failed")
			}
			if pr.Addr == net.Addr(nil) {
				return nil, fmt.Errorf("peer.Addr is nil")
			}

			addSlice := strings.Split(pr.Addr.String(), ":")
			if addSlice[0] == "[" {
				// 本机地址
				sourceIP = "localhost"
			}
			sourceIP = addSlice[0]
		}

		defer func(begin time.Time) {
			took := time.Since(begin)
			fmt.Println("sourceIp", sourceIP, name, "took", took)
		}(start)

		// 打印 request, response sourceIP
		resp, err := endpoint(ctx, request)
		if response, ok := resp.(*pb.LukeResponse); ok {
			response.RequestId = ksuid.New().String()
			if response.Retcode != pb.Retcode_OK {
				// 告警
				fmt.Println(response.Error)
			}
			resp = response
		}
		return resp, err
	}
}

var LabeleMiddlewareWithTracer = func(tracer otg.Tracer) LabeledMiddleware {
	return func(name string, in endpoint.Endpoint) endpoint.Endpoint {
		return gokitOt.TraceServer(tracer, name)(in)
	}
}
