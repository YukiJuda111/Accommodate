package main

import (
	"github.com/go-micro/plugins/v4/registry/consul"
	"postHouses/handler"
	pb "postHouses/proto"

	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	mysql "postHouses/model/mysql"
	redis "postHouses/model/redis"
)

var (
	service = "posthouses"
	version = "latest"
)

func main() {
	_, err := mysql.GormInit()
	if err != nil {
		return
	}
	redis.RedisInit()
	// Create service
	consulReg := consul.NewRegistry()
	srv := micro.NewService()
	srv.Init(
		micro.Name(service),
		micro.Version(version),
		micro.Registry(consulReg),
		micro.Address("127.0.0.1:5491"),
	)

	// Register handler
	if err := pb.RegisterPostHousesHandler(srv.Server(), new(handler.PostHouses)); err != nil {
		logger.Fatal(err)
	}
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
