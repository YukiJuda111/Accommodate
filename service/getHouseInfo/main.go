package main

import (
	"getHouseInfo/handler"
	pb "getHouseInfo/proto"
	"github.com/go-micro/plugins/v4/registry/consul"

	mysql "getHouseInfo/model/mysql"
	redis "getHouseInfo/model/redis"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
)

var (
	service = "gethouseinfo"
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
		micro.Address("127.0.0.1:5492"),
	)

	// Register handler
	if err := pb.RegisterGetHouseInfoHandler(srv.Server(), new(handler.GetHouseInfo)); err != nil {
		logger.Fatal(err)
	}
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
