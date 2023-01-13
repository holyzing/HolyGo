package luke

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/holyzing/HolyGo/go-kit/luke/endpoint"
	"github.com/holyzing/HolyGo/go-kit/luke/service"
	"github.com/holyzing/HolyGo/go-kit/luke/trace"
	"github.com/holyzing/HolyGo/go-kit/luke/transport"

	pb "github.com/holyzing/HolyGo/go-kit/luke/api/v1"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

const (
	maxMsgSize = 1024 * 1024 * 20
)

func InterruptHandler(errc chan<- error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	terminateError := fmt.Errorf("%s", <-c)

	// Place whatever shutdown handling you want here

	errc <- terminateError
}

func TestRun(t *testing.T) {
	errc := make(chan error)
	service := service.GetDefaultLukeService()
	endpoint := endpoint.NewLukeEndpointWithService(service)

	// go RunHttpServer(errc, endpoint)
	go RunMetricsServer(errc)
	go RunGrpcServer(errc, endpoint)
	go InterruptHandler(errc)
	// Run!
	<-errc
	println("--------------------> closing tracer ")
	trace.TracerCloser.Close()
	time.Sleep(5 * time.Second)
	println("--------------------> tracer closed !")
	log.Println("exit")
}

func RunHTTPServer(errc chan error, endpoints endpoint.LukeEndPoints) {
	log.Println("transport", "HTTP")
	h := transport.MakeHTTPHandler(endpoints)

	errc <- http.ListenAndServe("10.240.2.127:8080", h)
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

	// TODO 这是要干啥呢 ????
	// reflection.Register(s)

	ln, err := net.Listen("tcp", "10.240.2.127:9090")
	if err != nil {
		errc <- err
		return
	}
	errc <- s.Serve(ln)

	// 微服务注册: 服务名称 服务地址 有哪些方法 ?
	// 微服务发现: 服务名称 去获取服务地址
}

// prometheus metrics
func RunMetricsServer(errc chan error) {
	log.Println("transport", "metrics")
	h := http.NewServeMux()
	h.Handle("/metrics", promhttp.Handler())
	errc <- http.ListenAndServe("10.240.2.127:7070", h)
}

// H_Invincible_L

// truss

// httprouter
// mux
// m=http.NewServeMux
// m.Handle
// http.HandlerFunc
