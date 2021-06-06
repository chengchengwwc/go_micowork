package main

import (
	"context"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	consul2 "github.com/micro/go-plugins/registry/consul/v2"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"log"
	"user/common"
	server_user "user/proto/user"
)

func main() {
	// 注册中心
	consul := consul2.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})
	// 链路追踪
	t, io, err := common.NewTracer("go.micro.service.user", "localhost:6831")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)
	service := micro.NewService(
		micro.Name("go.micro.service.user.client"),
		micro.Version("latest"),
		micro.Address("127.0.0.1:8085"),
		//添加注册中心
		micro.Registry(consul),
		//绑定链路追踪
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
		// 添加限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(100)),
	)
	userService := server_user.NewUserService("go.micro.service.user", service.Client())
	response, err := userService.Login(context.TODO(), &server_user.UserLoginRequest{})
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println(response.IsSuccess)

}
