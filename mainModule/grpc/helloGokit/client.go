package helloGokit

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-kit/kit/transport/http"
)

// Direct: 直接调用服务端
// method:方法
// fullUrl: 完整的url http://localhost:8000
// enc: http.EncodeRequestFunc
// dec: http.DecodeResponseFunc
// requestStruct: 根据EndPoint定义的request结构体传参
func Direct(method string, fullUrl string, enc http.EncodeRequestFunc, dec http.DecodeResponseFunc, requestStruct interface{}) (interface{}, error) {
	// 1.解析url
	target, err := url.Parse(fullUrl)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// kit调用服务端拿到Client对象
	client := http.NewClient(strings.ToUpper(method), target, enc, dec)
	// 调用服务 client.Endpoint()返回一个可执行函数 传入context 和 请求数据结构体
	return client.Endpoint()(context.Background(), requestStruct)
}
