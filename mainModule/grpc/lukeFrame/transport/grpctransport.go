package transport

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log/logrus"
	"github.com/go-kit/kit/tracing/opentracing"
	"mainModule/grpc/lukeFrame/endpoint"
	"mainModule/grpc/lukeFrame/trace"
	"net/http"

	pb "mainModule/grpc/lukeFrame/proto"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/metadata"
)

type grpcServer struct {
	jobRead  grpctransport.Handler
	jobWrite grpctransport.Handler
}

func (le *grpcServer) JobWrite(ctx context.Context, req *pb.LukeRequest) (*pb.LukeResponse, error) {
	_, resp, err := le.jobRead.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.LukeResponse), err
}

func (le *grpcServer) JobRead(ctx context.Context, req *pb.LukeRequest) (*pb.LukeResponse, error) {
	_, resp, err := le.jobWrite.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.LukeResponse), err
}

// MakeGRPCServer makes a set of endpoints available as a gRPC LukeServiceServer.
func MakeGRPCServer(endpoints endpoint.LukeEndPoints) pb.LukeServiceServer {
	serverOptions := []grpctransport.ServerOption{
		grpctransport.ServerBefore(
			metadataToContext,
			opentracing.GRPCToContext(trace.OtTracer, "MakeGRPCServer", logrus.Logger{}),
		),
	}
	return &grpcServer{
		// lukeservice

		jobRead: grpctransport.NewServer(
			endpoints.JobReadEndPoint,
			DecodeLukeRequestFunc,
			EncodeLukeResponseFunc,
			serverOptions...,
		),
		jobWrite: grpctransport.NewServer(
			endpoints.JobWriteEndPoint,
			DecodeLukeRequestFunc,
			EncodeLukeResponseFunc,
			serverOptions...,
		),
	}
}

// NOTE 对于RPC 协议来说不需要进行 相应 Decode 和 Encode
// NOTE 但是对于 Http 协议 来说,需要将 heepRequest 转为 后端处理的 request, 并将 后端返回的 response 转为 http response
// NOTE 同理基于 thrift 协议的请求和响应也是如此, 其它协议也是如此

func DecodeLukeRequestFunc(_ context.Context, req interface{}) (request interface{}, err error) {
	request, ok := req.(*pb.LukeRequest)
	if ok {
		err = nil
	} else {
		err = fmt.Errorf("invalid luke request ")
	}
	return
}

// EncodeLukeResponseFunc is a transport/grpc.EncodeResponseFunc that converts a
// user-domain generic response to a gRPC generic reply. Primarily useful in a server.
func EncodeLukeResponseFunc(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.LukeResponse)
	return resp, nil
}

func metadataToContext(ctx context.Context, md metadata.MD) context.Context {
	for k, v := range md {
		if v != nil {
			// The key is added both in metadata format (k) which is all lower
			// and the http.CanonicalHeaderKey of the key so that it can be
			// accessed in either format

			// go-staticcheck ignore
			ctx = context.WithValue(ctx, k, v[0])
			ctx = context.WithValue(ctx, http.CanonicalHeaderKey(k), v[0])
		}
	}

	return ctx
}
