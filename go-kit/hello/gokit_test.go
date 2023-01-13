package hello

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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

// Direct: 直接调用服务端
// method:方法
// fullUrl: 完整的url http://localhost:8000
// enc: http.EncodeRequestFunc
// dec: http.DecodeResponseFunc
// requestStruct: 根据EndPoint定义的request结构体传参
func Direct(method, fullURL string, enc httpTransport.EncodeRequestFunc, dec httpTransport.DecodeResponseFunc, requestStruct interface{}) (interface{}, error) {
	// 1.解析url
	target, err := url.Parse(fullURL)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// kit调用服务端拿到Client对象
	client := httpTransport.NewClient(strings.ToUpper(method), target, enc, dec)
	// 调用服务 client.Endpoint()返回一个可执行函数 传入context 和 请求数据结构体
	return client.Endpoint()(context.Background(), requestStruct)
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
