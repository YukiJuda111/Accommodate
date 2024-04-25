package main

import (
	"RentHouse/web/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	// ------------gin框架的初始化----------------
	// 初始化路由
	router := gin.Default()
	// 静态资源处理
	router.Static("/home", "view")
	// 注册路由
	router.GET("/api/v1.0/session", controller.GetSession)
	router.GET("/api/v1.0/imagecode/:uuid", controller.GetImageCd)
	// 启动运行
	router.Run(":8080")
}
