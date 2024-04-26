package main

import (
	"github.com/go-micro/plugins/v4/registry/consul"
	"user/handler"
	"user/model"
	pb "user/proto"

	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
)

var (
	service = "user"
	version = "latest"
)

func main() {
	model.RedisInit()

	consulReg := consul.NewRegistry()
	// Create service
	srv := micro.NewService()
	srv.Init(
		micro.Name(service),
		micro.Version(version),
		micro.Registry(consulReg),
		micro.Address("127.0.0.1:5489"),
	)

	// Register handler
	if err := pb.RegisterUserHandler(srv.Server(), new(handler.User)); err != nil {
		logger.Fatal(err)
	}
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
