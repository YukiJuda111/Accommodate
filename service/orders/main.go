package main

import (
	"github.com/go-micro/plugins/v4/registry/consul"
	"orders/handler"
	pb "orders/proto"

	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"

	mysql "orders/model/mysql"
	redis "orders/model/redis"
)

var (
	service = "orders"
	version = "latest"
)

func main() {
	// Create service
	_, err := mysql.GormInit()
	if err != nil {
		return
	}
	redis.RedisInit()
	consulReg := consul.NewRegistry()

	srv := micro.NewService()
	srv.Init(
		micro.Name(service),
		micro.Version(version),
		micro.Registry(consulReg),
		micro.Address("127.0.0.1:5493"),
	)

	// Register handler
	if err := pb.RegisterOrdersHandler(srv.Server(), new(handler.Orders)); err != nil {
		logger.Fatal(err)
	}
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
