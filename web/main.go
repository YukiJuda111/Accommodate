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
	g := router.Group("/api/v1.0")
	{
		g.GET("/session", controller.GetSession)
		g.GET("/imagecode/:uuid", controller.GetImageCd)
		g.GET("/smscode/:phonenum", controller.GetSmscd)
	}
	// 启动运行
	router.Run(":8080")
}
