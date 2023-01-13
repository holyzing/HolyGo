package service

import (
	"context"

	pb "github.com/holyzing/HolyGo/go-kit/luke/api/v1"
)

type LukeService struct {
}

func (*LukeService) JobWrite(context.Context, *pb.LukeRequest) (*pb.LukeResponse, error) {
	var (
		err  error
		resp pb.LukeResponse
	)

	return &resp, err
}

func (LukeService) JobRead(context.Context, *pb.LukeRequest) (*pb.LukeResponse, error) {
	println("---- LukeService JobRead ----")
	var (
		err  error
		resp pb.LukeResponse
	)
	return &resp, err
}

var lukeService LukeService

// GetDefaultLukeService 接口类型是不能指针的, 而且接口类型的变量不能取地址, 接口本身是一个引用类型
func GetDefaultLukeService() pb.LukeServiceServer {
	return &lukeService
}

// TODO 接口的实现到底是用值传递还是指针传递 ????

func NewLukeService() pb.LukeServiceServer {
	var lukeService LukeService
	// var lukeService pb.LukeServiceServer = &LukeService{}
	return pb.LukeServiceServer(&lukeService)
}

// -------------------------------------------------------------------------------

// WrapService 可以对Service 中的 属性进一步封装, 注意属性可以是一个函数
func WrapService(lukeService pb.LukeServiceServer) pb.LukeServiceServer {
	return lukeService
}

// TODO 使用 GO-Kit 主要是为了实现  限流 融断 链路追踪
