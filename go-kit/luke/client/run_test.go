package client

import (
	"context"
	"log"
	"testing"
	"time"

	pb "github.com/holyzing/HolyGo/go-kit/luke/api/v1"
	"github.com/holyzing/HolyGo/go-kit/luke/client/api"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc"

	jaegercfg "github.com/uber/jaeger-client-go/config"
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
	defer client.tracerCloser.Close()
	ctx := context.Background()

	getJobInput := &api.GetJobInput{}

	/**
	// 链路追踪失败
	if res, ok := response.(endpoint.Failer); ok && res.Failed() != nil {
	*/

	if resp, err := client.GetJobEndpoint(ctx, getJobInput); err != nil {
		log.Fatalf("Server Error %v", err)
	} else {
		log.Printf("Server Response %v", resp)
	}
}

func TestTracer(t *testing.T) {
	// 创建一个连接配置
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "10.240.2.81:6831",
		},
		ServiceName: "TestTracer",
	}
	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	defer func() {
		closer.Close()
	}()
	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return
	}
	// 将创建的tracer设置为全局tracer，后续其他的任何地方都可以使用 opentracing.GlobalTracer() 返回的tracer代替
	opentracing.SetGlobalTracer(tracer)
	// 关闭一个tracer

	// ---------------------------------------------------------------------

	parentSpan := opentracing.StartSpan("测试延时时间")
	defer parentSpan.Finish()

	// opentracing.ChildOf(parentSpan.Context()) 用来说明父span
	subSpan := opentracing.StartSpan("测试20ms延时时间", opentracing.ChildOf(parentSpan.Context()))
	time.Sleep(time.Millisecond * 20)
	// 这里一定不可以使用 derfer来结束此span
	subSpan.Finish()

	// opentracing.ChildOf(parentSpan.Context()) 用来说明父span
	sub2Span := opentracing.StartSpan("测试30ms延时时间", opentracing.ChildOf(parentSpan.Context()))
	time.Sleep(time.Millisecond * 30)
	// 这里一定不可以使用 derfer来结束此span
	sub2Span.Finish()
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
