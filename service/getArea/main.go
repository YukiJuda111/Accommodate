package main

import (
	"getArea/handler"
	mysql "getArea/model/mysql"
	redis "getArea/model/redis"
	pb "getArea/proto"
	"github.com/go-micro/plugins/v4/registry/consul"

	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
)

var (
	service = "getarea"
	version = "latest"
)

func main() {
	_, err := mysql.GormInit()
	if err != nil {
		return
	}
	redis.RedisInit()

	consulReg := consul.NewRegistry()
	// Create service
	srv := micro.NewService()
	srv.Init(
		micro.Name(service),
		micro.Version(version),
		micro.Registry(consulReg),
		micro.Address("127.0.0.1:5490"),
	)

	// Register handler
	if err := pb.RegisterGetAreaHandler(srv.Server(), new(handler.GetArea)); err != nil {
		logger.Fatal(err)
	}
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
