package utils

import (
	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
)

func InitMicro() micro.Service {
	// 初始化consul配置
	consulReg := consul.NewRegistry()
	consulSev := micro.NewService(
		micro.Registry(consulReg),
	)
	return consulSev
}
