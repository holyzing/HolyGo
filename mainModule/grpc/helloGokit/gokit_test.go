package helloGokit

import (
	"fmt"
	"net/http"
	"testing"

	httpTransport "github.com/go-kit/kit/transport/http"
)

// Go kit 是一个微服务工具包集合。利用它提供的API和规范可以创建健壮、可维护性高的微服务体系

// 1、Service 这里就是我们的业务类、接口等相关信息存放
// 2、EndPoint 定义Request、Response格式，并可以使用装饰器(闭包)包装函数,以此来实现各个中间件嵌套
// 3、Transport 主要负责与HTTP、gRPC、thrift等相关逻辑

func TestGoKit(t *testing.T) {
	s := Server{}

	// 2.在用EndPoint/endpoint.go 创建业务服务
	hello := MakeServerEndPointHello(s)
	bye := MakeServerEndPointBye(s)

	// 3.使用 kit 创建 handler
	// 固定格式
	// 传入 业务服务 以及 定义的 加密解密方法

	// TODO 使用 Mux 扩展路由
	helloServer := httpTransport.NewServer(hello, HelloDecodeRequest, HelloEncodeResponse)
	sayServer := httpTransport.NewServer(bye, ByeDecodeRequest, ByeEncodeResponse)

	// 使用http包启动服务
	go http.ListenAndServe("0.0.0.0:8000", helloServer)

	go http.ListenAndServe("0.0.0.0:8001", sayServer)
	select {}
}

func TestClientDirect(t *testing.T) {
	i, err := Direct("GET", "http://127.0.0.1:8000",
		HelloEncodeRequestFunc,
		HelloDecodeResponseFunc,
		HelloRequest{Name: "songzhibin"},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	res, ok := i.(HelloResponse)
	if !ok {
		fmt.Println("no ok")
		return
	}
	fmt.Println(res)
}
