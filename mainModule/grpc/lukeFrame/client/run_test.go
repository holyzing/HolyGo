package client

import (
	"context"
	"log"
	"mainModule/grpc/lukeFrame/client/api"
	pb "mainModule/grpc/lukeFrame/proto"
	"testing"

	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc"
)

const (
	grpcAddress = "10.240.2.127:9090"
)

func TestCustomClient(t *testing.T) {

	conn, err := grpc.Dial(
		grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("can not connect and reason is \"%v\"", err)
	}
	defer conn.Close()
	// a := 'a'
	// TODO rune
	client := pb.NewLukeServiceClient(conn)
	req := &pb.LukeRequest{
		Body: &pb.LukeRequest_GetRequest{
			GetRequest: &pb.GetJobRequest{},
		},
	}

	resp, err := client.JobRead(context.Background(), req)
	if err != nil {
		log.Fatalf("Server Error %v", err)
	}
	log.Printf("Server Response %v", resp)
}

func TestLukeClient(t *testing.T) {
	option := Option{
		AccessToken: "AccessToken",
		Region:      "ALPHA",
		Env:         "DEVELOPMENT",
		Scheme:      "GRPC",
		Addr:        grpcAddress,
		Log:         nil,
	}
	client := New(option)
	ctx := context.Background()
	getJobInput := &api.GetJobInput{}
	if resp, err := client.GetJobEndpoint(ctx, getJobInput); err != nil {
		log.Fatalf("Server Error %v", err)
	} else {
		log.Printf("Server Response %v", resp)
	}
}

/**
真正要在企业内部使用的时候还须要一个注冊中心。
管理全部的服务。

初步计划使用 consul 存储数据。

由于consul 上面集成了许多的好东西。还有个简单的可视化的界面。


比etcd功能多些。

可是性能上面差一点。只是也很强悍了。
企业内部使用的注冊中心。已经足够了。
*/
