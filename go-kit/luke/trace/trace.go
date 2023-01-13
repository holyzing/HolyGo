package trace

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/uber/jaeger-client-go"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

var OtTracer opentracing.Tracer
var TracerCloser io.Closer

func init() {
	var addr string
	var err error
	// var closer io.Closer
	if os.Getenv("JAEGER_AGENT_HOST") != "" {
		addr = fmt.Sprintf("%s:5775", os.Getenv("JAEGER_AGENT_HOST"))
	} else {
		// addr = fmt.Sprintf("%s:%d", conf.Config().Server["jaeger"].Host, conf.Config().Server["jaeger"].Port)
		addr = fmt.Sprintf("%s:%d", "10.240.2.81", 6831)
	}
	cfg := config.Configuration{
		ServiceName: "Luke-Server",
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

	// NOTE TODO 全局变量,被内部使用 := 重新生成
	//OtTracer, closer, err := cfg.NewTracer(
	//	// config.Logger(log.Logger),
	//	//config.Metrics(metrics.NullFactory),
	//)

	tracer, closer, err := cfg.NewTracer(
		// config.Logger(log.Logger),
		//config.Metrics(metrics.NullFactory),
		config.Logger(jaeger.StdLogger),
	)
	if err != nil {
		log.Fatal(err)
	}
	TracerCloser = closer
	OtTracer = tracer
	opentracing.SetGlobalTracer(OtTracer)
	//defer closer.Close()
}
