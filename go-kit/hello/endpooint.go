package hello

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// endpoint.go 定义 Request、Response 格式, 并且可以使用闭包来实现各种中间件的嵌套
// 这里了解 protobuf 的比较好理解点
// 就是声明 接收数据和响应数据的结构体 并通过构造函数创建 在创建的过程当然可以使用闭包来进行一些你想要的操作啦

// HelloRequest 请求格式
type HelloRequest struct {
	Name string `json:"name"`
}

// HelloResponse 响应格式
type HelloResponse struct {
	Reply string `json:"reply"`
}

// ByeRequest 请求格式
type ByeRequest struct {
	Name string `json:"name"`
}

// ByeResponse 响应格式
type ByeResponse struct {
	Reply string `json:"reply"`
}

// Request 请求格式
type Request struct {
	Name string `json:"name"`
}

// Response 响应格式
type Response struct {
	Reply string `json:"reply"`
}

// MakeServerEndPointHello 这里创建 “构造函数“ hello方法的业务处理
func MakeServerEndPointHello(s IServer) endpoint.Endpoint {
	// 这里使用闭包,可以在这里做一些中间件业务的处理
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// request 是在对应请求来时传入的参数(这里的request 实际上是等下我们要将的Transport中一个decode函数中处理获得的参数)
		// 这里进行以下断言
		r, ok := request.(HelloRequest)
		if !ok {
			return Response{}, nil
		}
		// 这里实际上就是调用我们在Server/server.go中定义的业务逻辑
		// 我们拿到了 Request.Name 那么我们就可以调用我们的业务 Server.IServer 中的方法来处理这个数据并返回
		// 具体的业务逻辑具体定义....
		return HelloResponse{Reply: s.Hello(r.Name)}, nil
		// response 这里返回的response 可以返回任意的 不过根据规范是要返回我们刚才定义好的返回对象

	}
}

// MakeServerEndPointBye 这里创建构造函数 Bye方法的业务处理
func MakeServerEndPointBye(s IServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r, ok := request.(ByeRequest)
		if !ok {
			return Response{}, nil
		}
		return ByeResponse{Reply: s.Bye(r.Name)}, nil
	}
}

// 将当前的 Server 和 Request 以及 Response 都替换为 Rpc 生成的 Request和Response 即可
// 这样就需要在调用Endpoint 之前 将 RPCRequest转换为Endpoint Request, 业务处理之后，将 RPCResponse 转换为 Endpoint Response
