package endpoint

import (
	"context"
	"fmt"
	"reflect"

	"github.com/holyzing/HolyGo/go-kit/luke/trace"

	pb "github.com/holyzing/HolyGo/go-kit/luke/api/v1"
	"github.com/holyzing/HolyGo/go-kit/luke/middleware"

	"github.com/go-kit/kit/endpoint"
)

type LukeEndPoints struct {
	JobWriteEndPoint endpoint.Endpoint
	JobReadEndPoint  endpoint.Endpoint
}

func (le LukeEndPoints) JobWrite(ctx context.Context, req *pb.LukeRequest) (*pb.LukeResponse, error) {
	resp, err := le.JobWriteEndPoint(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.LukeResponse), err
}

func (le LukeEndPoints) JobRead(ctx context.Context, req *pb.LukeRequest) (*pb.LukeResponse, error) {
	println("---- LukeEndPoints JobRead ----")
	resp, err := le.JobReadEndPoint(ctx, req)
	if err != nil {
		return nil, err
	}
	// cannot convert (type interface {}) to type pb.LukeResponse: need type assertion
	// return pb.LukeResponse(resp), err
	return resp.(*pb.LukeResponse), err
}

// -------------------------------------------------------------------------------

func WrapJobReadServiceToEndPoint(ls pb.LukeServiceServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		println("---- WrapJobReadServiceToEndPoint ----")
		req := request.(*pb.LukeRequest)
		response, err = ls.JobRead(ctx, req)
		if err != nil {
			return nil, err
		}
		return
	}
}

func WrapJobWriteServiceToEndPoint(ls pb.LukeServiceServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		println("---- WrapJobWriteServiceToEndPoint ----")
		req := request.(*pb.LukeRequest)
		resp, err := ls.JobWrite(ctx, req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}

// -------------------------------------------------------------------------------

// NOTE 封装方式 可以分为单个外部封装, 也可以统一内部封装
func (le *LukeEndPoints) WrapEndpointsWithMiddleware(middleware endpoint.Middleware) {
	// leType := reflect.TypeOf(le).Elem()
	leValue := reflect.ValueOf(le).Elem()
	epType := reflect.TypeOf(endpoint.Endpoint(nil))
	epKind := epType.Kind()

	for i := 0; i < leValue.NumField(); i++ {
		// TODO 使用反射去判断, 包装
		var field = leValue.Field(i)
		if field.Kind() == epKind {
			ep := middleware(field.Interface().(endpoint.Endpoint))
			field.Set(reflect.ValueOf(ep))
			// leValue.FieldByName(field.Name).Set(reflect.ValueOf(ep))
		}
	}

	/*
		for i := 0; i < leType.NumField(); i++ {
			// TODO 使用反射去判断, 包装
			var field = leType.Field(i)
			if field.Type.Kind() == epKind {
				ep := middleware(leValue.Index(i).Interface().(endpoint.Endpoint))
				leValue.FieldByName(field.Name).Set(reflect.ValueOf(ep))
			}
		}
	*/
}

func (le *LukeEndPoints) WrapEndpointsWithLabelMiddleware(middleware middleware.LabeledMiddleware, excluded ...string) {
	included := map[string]struct{}{
		"JobRead":  struct{}{},
		"JobWrite": struct{}{},
	}

	for _, ex := range excluded {
		if _, ok := included[ex]; !ok {
			panic(fmt.Sprintf("Excluded endpoint '%s' does not exist; see middlewares/endpoints.go", ex))
		}
		delete(included, ex)
	}

	for inc := range included {
		if inc == "JobRead" {
			le.JobReadEndPoint = middleware("JobRead", le.JobReadEndPoint)
		} else if inc == "JobWrite" {
			le.JobWriteEndPoint = middleware("JobWrite", le.JobWriteEndPoint)
		}
	}
}

func NewLukeEndpointWithService(service pb.LukeServiceServer) LukeEndPoints {
	var lukeendpoint = LukeEndPoints{
		JobReadEndPoint:  WrapJobReadServiceToEndPoint(service),
		JobWriteEndPoint: WrapJobWriteServiceToEndPoint(service),
	}
	lukeendpoint.WrapEndpointsWithMiddleware(middleware.Auth0Middleware)
	lukeendpoint.WrapEndpointsWithLabelMiddleware(middleware.LoggingMiddleware)

	lukeendpoint.JobReadEndPoint = middleware.LabeleMiddlewareWithTracer(trace.OtTracer)("JobReadEndPoint", lukeendpoint.JobReadEndPoint)
	lukeendpoint.JobWriteEndPoint = middleware.LabeleMiddlewareWithTracer(trace.OtTracer)("JobWriteEndPoint", lukeendpoint.JobWriteEndPoint)

	// in.WrapAllExcept(opentracing.TraceServer(trace.OtTracer, "Generic"), "Admin")
	return lukeendpoint
}
