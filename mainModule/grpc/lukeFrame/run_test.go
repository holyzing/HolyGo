package lukeFrame

import (
	"log"
	"mainModule/grpc/lukeFrame/endpoint"
	"mainModule/grpc/lukeFrame/service"
	"mainModule/grpc/lukeFrame/transport"
	"net"
	"net/http"
	"testing"

	pb "mainModule/grpc/lukeFrame/proto"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

const (
	maxMsgSize = 1024 * 1024 * 20
)

func TestRun(t *testing.T) {
	errc := make(chan error)
	service := service.GetDefaultLukeService()
	endpoint := endpoint.NewLukeEndpointWithService(service)

	go RunHttpServer(errc, endpoint)

	// Run!
	log.Println("exit", <-errc)
}

func RunHttpServer(errc chan error, endpoints endpoint.LukeEndPoints) {
	log.Println("transport", "HTTP")
	h := transport.MakeHTTPHandler(endpoints)

	errc <- http.ListenAndServe("10.240.2.137:8080", h)
}

func RunGrpcServer(errc chan error, endpoints endpoint.LukeEndPoints) {
	log.Println("transport", "gRPC")
	var opts = []grpc.ServerOption{
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	}

	s := grpc.NewServer(opts...)
	grpc_prometheus.Register(s)
	srv := transport.MakeGRPCServer(endpoints)
	pb.RegisterLukeServiceServer(s, srv)

	ln, err := net.Listen("tcp", "10.240.2.137:9090")
	if err != nil {
		errc <- err
		return
	}
	errc <- s.Serve(ln)
}

// prometheus metrics
func RunMetricsServer(errc chan error) {
	log.Println("transport", "metrics")
	h := http.NewServeMux()
	h.Handle("/metrics", promhttp.Handler())
	errc <- http.ListenAndServe("10.240.2.137:7070", h)
}
