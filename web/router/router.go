package router

import (
	"RentHouse/web/controller"
	"RentHouse/web/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init() {
	// 初始化路由
	router := gin.Default()
	// 初始化session
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	// 静态资源处理
	router.Static("/home", "view")
	// 注册路由
	g := router.Group("/api/v1.0")
	{
		g.GET("/imagecode/:uuid", controller.GetImageCd)
		g.GET("/smscode/:phonenum", controller.GetSmscd)
		g.POST("/users", controller.PostRet)
		g.GET("/areas", controller.GetArea)
		g.POST("/sessions", controller.PostLogin)

		// 使用中间件,之后的接口都不需要校验Session了
		g.Use(middleware.LoginFilter())
		g.GET("/session", controller.GetSession)
		g.DELETE("/session", controller.DeleteSession)
		g.GET("/user", controller.GetUserInfo)
		g.PUT("/user/name", controller.PutUserInfo)
		g.POST("/user/avatar", controller.PostAvatar)
		g.POST("/user/auth", controller.PutUserAuth)
		g.GET("/user/auth", controller.GetUserInfo)
		g.GET("/user/houses", controller.GetHouses)
	}
	// 启动运行
	router.Run(":8080")
}
