package main

import (
	"getCaptcha/handler"
	pb "getCaptcha/proto"
	"github.com/go-micro/plugins/v4/registry/consul"

	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
)

var (
	service = "getcaptcha"
	version = "latest"
)

func main() {
	consulReg := consul.NewRegistry()

	// Create service
	srv := micro.NewService()
	srv.Init(
		micro.Name(service),
		micro.Version(version),
		micro.Registry(consulReg),
		micro.Address("127.0.0.1:5488"),
	)

	// Register handler
	if err := pb.RegisterGetCaptchaHandler(srv.Server(), new(handler.GetCaptcha)); err != nil {
		logger.Fatal(err)
	}
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
