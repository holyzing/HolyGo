package transport

import (
	"context"
	"fmt"
	"mainModule/grpc/lukeFrame/client/api"
	pb "mainModule/grpc/lukeFrame/proto" // 客户端服务断都需要

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

func NewGRPCClient(conn *grpc.ClientConn, options ...grpctransport.ClientRequestFunc) Set {
	var clientOptions = []grpctransport.ClientOption{}

	for _, option := range options {
		clientOptions = append(clientOptions, grpctransport.ClientBefore(option))
	}

	// Each individual endpoint is an grpc/transport.Client (which implements
	// endpoint.Endpoint) that gets wrapped with various middlewares. If you
	// made your own client library, you'd do this work there, so your server
	// could rely on a consistent set of client behavior.
	var createJobEndpoint endpoint.Endpoint
	{
		createJobEndpoint = grpctransport.NewClient(
			conn,
			"luke.LukeService",
			"JobWrite",
			encodeGRPCCreateJobRequest,
			decodeGRPCCreateJobResponse,
			pb.LukeResponse{},
			clientOptions...,
		).Endpoint()
	}

	var getJobEndpoint endpoint.Endpoint
	{
		getJobEndpoint = grpctransport.NewClient(
			conn,
			"luke.LukeService",
			"JobRead",
			encodeGRPCGetJobRequest,
			decodeGRPCGetJobResponse,
			pb.LukeResponse{},
			clientOptions...,
		).Endpoint()
	}

	return Set{
		CreateJobEndpoint: createJobEndpoint,
		GetJobEndpoint:    getJobEndpoint,
	}
}

// encodeGRPCGetJobRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain concat request to a gRPC concat request. Primarily useful in a
// client.
func encodeGRPCGetJobRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*api.GetJobInput)
	return &pb.LukeRequest{
		Method:       "GetJob",
		User:         req.User,
		Organization: req.Organization,
		TenantName:   req.TenantName,
		Body: &pb.LukeRequest_GetRequest{
			GetRequest: &pb.GetJobRequest{
				Id:       req.ID,
				Handle:   req.Handle,
				Fields:   req.Fields,
				Combined: req.Combined,
			},
		},
	}, nil
}

// decodeGRPCGetJobResponse is a transport/grpc.DecodeResponseFunc that converts
// a gRPC concat reply to a user-domain concat response. Primarily useful in a
// client.
func decodeGRPCGetJobResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.LukeResponse)
	output := &api.GetJobOutput{
		RequestID: reply.RequestId,
	}
	if reply.Retcode != pb.Retcode_OK && reply.Error != nil {
		return output, fmt.Errorf("%s:%s", reply.Error.Code, reply.Error.Message)
	}
	output.Body = convertJobDetail(reply.GetJobDetails())

	return output, nil
}

func convertJobDetail(detail *pb.JobDetailResponse) *api.JobDetailResponse {
	var jobDetails = make([]*api.CreateJobInput, 0)

	for _, d := range detail.GetJobDetails() {
		jobDetail := api.CreateJobInput{
			User:       d.User,
			Sync:       d.Sync,
			Cores:      d.Cores,
			SysPrio:    d.SysPrio,
			Entrypoint: d.Entrypoint,
			Options:    d.Options,
		}
		jobDetails = append(jobDetails, &jobDetail)
	}

	return &api.JobDetailResponse{
		JobDetails: jobDetails,
	}
}

// encodeGRPCCreateJobRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain sum request to a gRPC sum request. Primarily useful in a client.
func encodeGRPCCreateJobRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*api.CreateJobInput)
	return &pb.LukeRequest{
		Method:       "CreateJob",
		User:         req.User,
		Organization: req.Organization,
		TenantName:   req.TenantName,
		Body: &pb.LukeRequest_CreateRequest{
			CreateRequest: &pb.CreateJobRequest{
				Sync:       req.Sync,
				Cores:      req.Cores,
				SysPrio:    req.SysPrio,
				Entrypoint: req.Entrypoint,
				Options:    req.Options,
			},
		},
	}, nil
}

// decodeGRPCCreateJobResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC sum reply to a user-domain sum response. Primarily useful in a client.
func decodeGRPCCreateJobResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.LukeResponse)
	output := &api.CreateJobOutput{
		RequestID: reply.RequestId,
	}
	if reply.Retcode != pb.Retcode_OK && reply.Error != nil {
		return output, fmt.Errorf("%s:%s", reply.Error.Code, reply.Error.Message)
	}
	output.Body = &api.JobInfoResponse{
		JobID:     reply.GetJobInfo().GetJobId(),
		JobHandle: reply.GetJobInfo().GetJobHandle(),
	}

	return output, nil
}
